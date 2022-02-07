package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

func main() {
	log.Printf("Starting.")

	stepDuration := time.Duration(10 * time.Second)

	createClientCtx, cancelCreateClient := context.WithTimeout(context.Background(), stepDuration)
	defer cancelCreateClient()
	client, err := pubsub.NewClient(createClientCtx, "testproject")
	if err != nil {
		log.Fatal(err)
	}

	// Create topic and subscription to it
	topic := client.Topic("topic1")
	createSubCtx, cancelCreateSub := context.WithTimeout(context.Background(), stepDuration)
	defer cancelCreateSub()
	sub, err := client.CreateSubscription(createSubCtx, "subscription1", pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Publish
	publishCtx, cancelPublish := context.WithTimeout(context.Background(), stepDuration)
	defer cancelPublish()
	res := topic.Publish(publishCtx, &pubsub.Message{
		Data: []byte("hello world"),
	})

	// Get publish results
	getPublishResCtx, cancelGetPublishRes := context.WithTimeout(context.Background(), stepDuration)
	defer cancelGetPublishRes()
	msgID, err := res.Get(getPublishResCtx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Published message with ID: %s", msgID)

	// Use a callback to receive messages via subscription1.
	receiveCtx, cancelReceive := context.WithTimeout(context.Background(), stepDuration)
	defer cancelReceive()
	err = sub.Receive(receiveCtx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Received message: %v", m)
		m.Ack() // Acknowledge that we've consumed the message.
		log.Printf("Acked message.")
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Program completed.")
}
