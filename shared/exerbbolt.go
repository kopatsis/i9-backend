package shared

import (
	"context"
	"encoding/json"
	"log"

	bolt "go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bucketName = "CacheBucket"

func GetExersHelper(database *mongo.Database, boltDB *bolt.DB) ([]Exercise, error) {
	var exercises []Exercise

	err := boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		v := b.Get([]byte("Exercise"))
		if v != nil {
			err := json.Unmarshal(v, &exercises)
			if err == nil {
				return nil
			}
			log.Printf("Failed to unmarshal from bbolt: %v, fetching from MongoDB", err)
		}

		findOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
		cursor, err := database.Collection("exercise").Find(context.Background(), bson.D{}, findOptions)
		if err != nil {
			return err
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &exercises); err != nil {
			return err
		}

		data, err := json.Marshal(exercises)
		if err != nil {
			return err
		}

		err = b.Put([]byte("Exercise"), data)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return exercises, nil
}
