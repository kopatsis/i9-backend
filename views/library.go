package views

import (
	"context"
	"fulli9/shared"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetLibrary(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		var wg sync.WaitGroup

		errChan := make(chan error, 3)
		var errGroup *multierror.Error

		user, exers, strs := shared.User{}, []shared.Exercise{}, []shared.Stretch{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			user, err = getUserHelper(database, userID)
			if err != nil {
				errChan <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			exers, err = shared.GetExersHelper(database, boltDB)
			if err != nil {
				errChan <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			strs, err = shared.GetStretchHelper(database, boltDB)
			if err != nil {
				errChan <- err
			}
		}()

		wg.Wait()
		close(errChan)

		hasErr := false
		for err := range errChan {
			if err != nil {
				errGroup = multierror.Append(errGroup, err)
				hasErr = true
			}
		}

		if hasErr {
			c.JSON(400, gin.H{
				"Error": "Issue with querying",
				"Exact": errGroup.Error(),
			})
			return
		}

		retExer, retStr := []shared.RetLibraryExer{}, []shared.RetLibraryStr{}

		for _, exer := range exers {
			blocked := slices.Contains(user.BannedExercises, exer.ID.Hex())
			var fav float32
			if val, ok := user.ExerFavoriteRates[exer.ID.Hex()]; ok {
				fav = val
			} else {
				fav = 1
			}

			retExer = append(retExer, shared.RetLibraryExer{
				ID:         exer.ID.Hex(),
				Name:       exer.Name,
				Parent:     exer.Parent,
				Blocked:    blocked,
				Favoritism: fav,
				BodyParts:  exer.BodyParts,
			})
		}

		for _, str := range strs {
			blocked := slices.Contains(user.BannedStretches, str.ID.Hex())

			retStr = append(retStr, shared.RetLibraryStr{
				ID:        str.ID.Hex(),
				Name:      str.Name,
				Type:      str.Status,
				Blocked:   blocked,
				BodyParts: str.BodyParts,
			})
		}

		c.JSON(200, gin.H{
			"Exers": &retExer,
			"Strs":  &retStr,
		})

	}
}

func getUserHelper(database *mongo.Database, userID string) (shared.User, error) {
	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return shared.User{}, err
	}

	collection := database.Collection("user")
	filter := bson.D{{Key: "_id", Value: id}}
	var user shared.User

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return shared.User{}, err
	}

	return user, nil
}
