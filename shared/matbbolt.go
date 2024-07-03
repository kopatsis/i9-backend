package shared

import (
	"context"
	"encoding/json"
	"log"

	bolt "go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMatrixHelper(database *mongo.Database, boltDB *bolt.DB) (TypeMatrix, error) {
	var matrix TypeMatrix

	err := boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		v := b.Get([]byte("Matrix"))
		if v != nil {
			err := json.Unmarshal(v, &matrix)
			if err == nil {
				return nil
			}
			log.Printf("Failed to unmarshal from bbolt: %v, fetching from MongoDB", err)
		}

		cursor, err := database.Collection("typematrix").Find(context.Background(), bson.D{})
		if err != nil {
			return err
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &matrix); err != nil {
			return err
		}

		data, err := json.Marshal(matrix)
		if err != nil {
			return err
		}

		err = b.Put([]byte("Matrix"), data)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return TypeMatrix{}, err
	}

	return matrix, nil
}
