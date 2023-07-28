package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alancesar/downloader/pkg/media"
	"github.com/alancesar/downloader/pkg/redgifs"
	"github.com/alancesar/downloader/pkg/ticker"
	"github.com/alancesar/downloader/pkg/transport"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalln("failed to open amqp connection:", err)
	}

	defer func() {
		_ = amqpConnection.Close()
	}()

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

	producerChannel, err := amqpConnection.Channel()
	if err != nil {
		return
	}

	defer func() {
		_ = producerChannel.Close()
	}()

	consumer := func(m media.Media) error {
		gif, err := redGIFsClient.GetGIFByURL(m.URL)
		if err != nil {
			return fmt.Errorf("on retrieve gif data: %w", err)
		}

		gifMedia := gif.ToMedia()
		body, err := json.Marshal(gifMedia)
		if err != nil {
			return fmt.Errorf("on unmarmal media: %w", err)
		}

		if err := producerChannel.PublishWithContext(
			ctx,
			"media",
			"downloads",
			false,
			false,
			amqp.Publishing{
				Headers: map[string]interface{}{
					"provider": "redgifs",
				},
				ContentType: "application/json",
				Body:        body,
			},
		); err != nil {
			return fmt.Errorf("on publish gif media: %w", err)
		}

		return nil
	}

	consumerChannel, err := amqpConnection.Channel()
	if err != nil {
		log.Fatalln("failed to open amqp channel:", err)
	}

	defer func() {
		_ = consumerChannel.Close()
	}()

	if err := consumerChannel.Qos(200, 0, false); err != nil {
		log.Fatalln("failed to set channel qos:", err)
	}

	messages, err := consumerChannel.Consume(
		"media.redgifs",
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

	for message := range messages {
		var m media.Media
		if err := json.Unmarshal(message.Body, &m); err != nil {
			log.Println("failed to unmarshal message", err)
			_ = message.Ack(false)
			continue
		}

		if err := consumer(m); err != nil {
			fmt.Println(err)
			_ = message.Nack(false, true)
			continue
		}

		_ = message.Ack(false)
	}
}
