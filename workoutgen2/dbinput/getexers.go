package dbinput

import (
	"fulli9/shared"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetExersDB(database *mongo.Database, boltDB *bbolt.DB) (map[string]shared.Exercise, error) {

	exers, err := shared.GetExersHelper(database, boltDB)
	if err != nil {
		return nil, err
	}

	allexer := map[string]shared.Exercise{}
	for _, exer := range exers {
		allexer[exer.ID.Hex()] = exer
	}

	return allexer, nil
}
