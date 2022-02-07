package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	log.Printf("Starting.")

	client, err := pubsub.NewClient(context.Background(), "testproject")
	if err != nil {
		log.Fatal(err)
	}

	// Create topic and subscription to it
	topic, err := client.CreateTopic(context.Background(), "topic1")
	if err != nil {
		log.Fatal(err)
	}
	sub, err := client.CreateSubscription(context.Background(), "subscription1", pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Publish
	res := topic.Publish(context.Background(), &pubsub.Message{
		Data: []byte("hello world"),
	})

	// Get publish results
	msgID, err := res.Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Published message with ID: %s", msgID)

	for3Seconds, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Use a callback to receive messages via subscription1.
	err = sub.Receive(for3Seconds, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Received message: %v", m)
		m.Ack() // Acknowledge that we've consumed the message.
		log.Printf("Acked message.")
		cancel()
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Program completed.")
}
