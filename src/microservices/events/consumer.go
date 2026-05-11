package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func startMovieConsumer() {
	startConsumer("movie-events", "movie-event-consumer")
}

func startUserConsumer() {
	startConsumer("user-events", "user-event-consumer")
}

func startPaymentConsumer() {
	startConsumer("payment-events", "payment-event-consumer")
}

func startConsumer(topic, groupId string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    topic,
		GroupID:  groupId,
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
