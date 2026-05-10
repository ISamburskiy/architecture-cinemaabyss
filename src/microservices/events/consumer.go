package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func startMovieConsumer() {
	startConsumer("movie-events")
}

func startUserConsumer() {
	startConsumer("user-events")
}

func startPaymentConsumer() {
	startConsumer("payment-events")
}

func startConsumer(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    topic,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading from %s: %v", topic, err)
			continue
		}

		var prettyPayload map[string]interface{}
	if err := json.Unmarshal(m.Value, &prettyPayload); err == nil {
			prettyJSON, _ := json.MarshalIndent(prettyPayload, "  ", "  ")
			log.Printf("CONSUMED [%s] - Offset: %d\n%s",
			topic, m.Offset, string(prettyJSON))
	} else {
			// Если не JSON, выводим как строку
			log.Printf("CONSUMED [%s] - Offset: %d, Raw: %s",
		topic, m.Offset, string(m.Value))
	}
	}
}
