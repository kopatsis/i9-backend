package views

import (
	"context"
	"fmt"
	"fulli9/shared"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetHistory(database *mongo.Database) gin.HandlerFunc {
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

		errChan := make(chan error, 4)
		var errGroup *multierror.Error

		wos, strwos, ex, st := []shared.Workout{}, []shared.StretchWorkout{}, map[string]string{}, map[string]string{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			wos, err = getWOsHelper(database, userID)
			if err != nil {
				errChan <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			strwos, err = getStrWOsHelper(database, userID)
			if err != nil {
				errChan <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			ex, err = getNamesHelper(database, "exercise")
			if err != nil {
				errChan <- err
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			st, err = getNamesHelper(database, "stretch")
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

		for i, strwo := range strwos {
			for j, static := range strwo.Statics {
				strwo.Statics[j] = st[static]
			}
			for j, dynamic := range strwo.Dynamics {
				strwo.Dynamics[j] = st[dynamic]
			}
			strwos[i] = strwo
		}

		for i, wo := range wos {
			for j, static := range wo.Statics {
				wo.Statics[j] = st[static]
			}
			for j, dynamic := range wo.Dynamics {
				wo.Dynamics[j] = st[dynamic]
			}
			for j, round := range wo.Exercises {
				for k, id := range round.ExerciseIDs {
					wo.Exercises[j].ExerciseIDs[k] = ex[id]
				}
			}
			wos[i] = wo
		}

		c.JSON(200, gin.H{
			"Stretch": &strwos,
			"Workout": &wos,
		})

	}
}

func getWOsHelper(database *mongo.Database, userID string) ([]shared.Workout, error) {
	collection := database.Collection("workout")

	optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

	filterWO := bson.D{
		{Key: "userid", Value: userID},
	}

	filterWO = append(filterWO, bson.E{Key: "status", Value: bson.D{{Key: "$ne", Value: "Archived"}}})

	cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var pastWorkouts []shared.Workout
	err = cursor.All(context.Background(), &pastWorkouts)
	if err != nil {
		return nil, err
	}

	return pastWorkouts, nil
}

func getStrWOsHelper(database *mongo.Database, userID string) ([]shared.StretchWorkout, error) {
	collection := database.Collection("stretchworkout")

	optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

	filterWO := bson.D{
		{Key: "userid", Value: userID},
	}

	filterWO = append(filterWO, bson.E{Key: "status", Value: bson.D{{Key: "$ne", Value: "Archived"}}})

	cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var pastWorkouts []shared.StretchWorkout
	err = cursor.All(context.Background(), &pastWorkouts)
	if err != nil {
		return nil, err
	}

	return pastWorkouts, nil
}

func getNamesHelper(database *mongo.Database, collectionName string) (map[string]string, error) {
	idNameMap := make(map[string]string)

	collection := database.Collection(collectionName)
	projection := bson.D{{Key: "_id", Value: 1}, {Key: "name", Value: 1}}
	opts := options.Find().SetProjection(projection)

	cur, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		id := result["_id"].(primitive.ObjectID).Hex()
		name := result["name"].(string)
		idNameMap[id] = name
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return idNameMap, nil
}
