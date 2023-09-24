package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type WebHookPayload struct {
	Url       string `json:"url"`
	WebhookId string `json:"webhookId"`
	Data      struct {
		Id      string `json:"id"`
		Payment string `json:"payment"`
		Event   string `json:"event"`
		Date    string `json:"created"`
	} `json:"data"`
}

func Subscribe(ctx context.Context, client *redis.Client, webhookQueue chan WebHookPayload) error {
	pubsub := client.Subscribe(ctx, "payments")

	defer func(pubsub *redis.PubSub) {
		if err := pubsub.Close(); err != nil {
			log.Println("Error closing pubsub: ", err)
		}
	}(pubsub)

	var payload WebHookPayload

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(msg.Payload), &payload)

		if err != nil {
			log.Println("Error unmarshalling JSON: ", err)
			continue
		}

		webhookQueue <- payload
	}
}
