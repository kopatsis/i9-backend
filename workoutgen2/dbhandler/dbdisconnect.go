package dbhandler

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func DisConnectDB(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
