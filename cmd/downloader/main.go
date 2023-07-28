package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alancesar/downloader/internal/database"
	"github.com/alancesar/downloader/internal/storage"
	"github.com/alancesar/downloader/pkg/media"
	"github.com/alancesar/downloader/pkg/redgifs"
	"github.com/alancesar/downloader/pkg/ticker"
	"github.com/alancesar/downloader/pkg/transport"
	"github.com/alancesar/downloader/usecase"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/schollz/progressbar/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dsn := filepath.Join(os.Getenv("STORAGE_ROOT"), "sqlite.db")
	sqliteConn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln("failed to open sqlite connection:", err)
	}
	gormDatabase, err := database.NewGorm(sqliteConn)
	if err != nil {
		log.Fatalln("failed to create sqlite database:", err)
	}

	localStorage := storage.NewLocalStorage(os.Getenv("STORAGE_ROOT"))
	progressBar := func(res *http.Response, description string) io.Writer {
		return progressbar.DefaultBytes(res.ContentLength, description)
	}

	defaultClient := &http.Client{
		Transport: transport.NewUserAgentRoundTripper("downloaddit/v0", http.DefaultTransport),
	}

	redGIFsAuthProvider := redgifs.NewClient(defaultClient)
	tokenTicker := ticker.NewToken(func() string {
		token, _ := redGIFsAuthProvider.RetrieveToken()
		return "Bearer " + token
	}, time.Minute*10)
	redGIFsAuthClient := &http.Client{
		Transport: transport.NewAuthorizationRoundTripper(tokenTicker.Get, defaultClient.Transport),
	}
	redGIFsClient := redgifs.NewClient(redGIFsAuthClient)
	interceptor := redgifs.NewInterceptor(redGIFsClient)
	interceptors := map[string]usecase.Interceptor{
		"redgifs": interceptor,
	}

	localDownloader := media.NewDownloader(localStorage, progressBar, defaultClient)
	useCase := usecase.NewDownload(localDownloader, gormDatabase, interceptors)

	consumer := func(m media.Media, provider string) error {
		if err := useCase.Execute(ctx, m, provider); err != nil {
			log.Println("failed to consume message:", err)
			if err := gormDatabase.SaveMedia(ctx, m); err == nil {
				log.Println("saving in media database")
				return nil
			}
			return err
		}

		return nil
	}

	for {
		amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err != nil {
			log.Fatalln("failed to open amqp connection:", err)
		}

		notify := amqpConnection.NotifyClose(make(chan *amqp.Error))

		channel, err := amqpConnection.Channel()
		if err != nil {
			log.Fatalln("failed to open amqp channel:", err)
		}

		if err := channel.Qos(10, 0, false); err != nil {
			log.Fatalln("failed to set channel qos:", err)
		}

		messages, err := channel.Consume(
			"media.downloads",
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalln("failed to consume media downloads queue:", err)
		}

		for {
			select {
			case err = <-notify:
				log.Println("connection lost", err)
				break
			case <-ctx.Done():
				log.Println("shutting down...")
				stop()
				_ = channel.Close()
				_ = amqpConnection.Close()
				log.Fatalln("good bye")
			case message := <-messages:
				var provider string
				if p := message.Headers["provider"]; p != nil {
					provider = p.(string)
				}

				var m media.Media
				if err := json.Unmarshal(message.Body, &m); err != nil {
					log.Println("failed to unmarshal message", err)
					_ = message.Ack(false)
					continue
				}

				if err := consumer(m, provider); err != nil {
					fmt.Println(err)
					_ = message.Nack(false, true)
					continue
				}

				_ = message.Ack(false)
			}
		}
	}
}
