package main

import (
	"context"
	"log"
	"os"
	
	"github.com/segmentio/kafka-go"
)

func produceToTopic(topic string, payload []byte) error {
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		os.Getenv("KAFKA_BROKERS"),
		topic,
		0,
	)
	if err != nil {
		log.Printf("Failed to dial leader for topic %s: %v", topic, err)
		return err
	}
	defer conn.Close()

	_, err = conn.WriteMessages(kafka.Message{Value: payload})
	if err != nil {
		log.Printf("Failed to write message to topic %s: %v", topic, err)
	} else {
		log.Printf("Produced to %s: %s", topic, string(payload))
	}
	return err
}
