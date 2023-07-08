package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alancesar/downloader/internal/database"
	"github.com/alancesar/downloader/internal/storage"
	"github.com/alancesar/downloader/pkg/media"
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

	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalln("failed to start amqp connection:", err)
	}

	channel, err := amqpConnection.Channel()
	if err != nil {
		log.Fatalln("failed to start amqp channel:", err)
	}

	dsn := filepath.Join(os.Getenv("STORAGE_ROOT"), "sqlite.db")
	sqliteConn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	gormDatabase, err := database.NewGorm(sqliteConn)

	localStorage := storage.NewLocalStorage(os.Getenv("STORAGE_ROOT"))
	progressBar := func(res *http.Response, description string) io.Writer {
		return progressbar.DefaultBytes(res.ContentLength, description)
	}

	defaultClient := &http.Client{
		Transport: transport.NewUserAgentRoundTripper("downloaddit/v0", http.DefaultTransport),
	}

	localDownloader := media.NewDownloader(localStorage, progressBar, defaultClient)
	useCase := usecase.NewDownload(localDownloader, gormDatabase)

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

	go func() {
		for message := range messages {
			var m media.Media
			if err := json.Unmarshal(message.Body, &m); err != nil {
				log.Println("failed to unmarshal message", err)
				_ = message.Ack(false)
				continue
			}

			if err := useCase.Execute(ctx, m); err != nil {
				log.Println("failed to consume message:", err)
				_ = message.Nack(false, true)
				continue
			}

			_ = message.Ack(false)
		}
	}()

	fmt.Println("all systems go!")

	<-ctx.Done()
	stop()

	fmt.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = amqpConnection.Close()
	fmt.Println("good bye")
}
