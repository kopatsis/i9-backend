package dbinput

import (
	"fulli9/shared"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStretchesDB(database *mongo.Database, boltDB *bbolt.DB) (map[string][]shared.Stretch, error) {

	allstretches := map[string][]shared.Stretch{}
	dynamics := []shared.Stretch{}
	statics := []shared.Stretch{}

	strs, err := shared.GetStretchHelper(database, boltDB)
	if err != nil {
		return nil, err
	}

	for _, str := range strs {
		if str.Status == "Dynamic" {
			dynamics = append(dynamics, str)
		} else {
			statics = append(statics, str)
		}
	}

	allstretches["Dynamic"] = dynamics
	allstretches["Static"] = statics

	return allstretches, nil
}
