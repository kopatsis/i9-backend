package workoutgen2

import (
	"context"
	"errors"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostWorkout(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.WorkoutRoute
		if err := c.ShouldBindJSON(&woHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}

		workout, err := WorkoutGen(woHandler.Time, woHandler.Difficulty, userID, database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}
		workoutRet := workout.(shared.Workout)

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workoutRet)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestWorkout(workoutRet, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workoutRet,
				"positions": resp,
			})
		}

	}
}

func PostStretchWorkout(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.StrWorkoutRoute
		if err := c.ShouldBindJSON(&woHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with stretch workout generator",
				"Exact": err.Error(),
			})
			return
		}

		workout, err := WorkoutGen(woHandler.Time, 0, userID, database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with stretch workout generator",
				"Exact": err.Error(),
			})
			return
		}
		workoutRet := workout.(shared.StretchWorkout)

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workoutRet)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestStrWorkout(workoutRet, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workoutRet,
				"positions": resp,
			})
		}
	}
}

func PostWorkoutRetry(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.WorkoutRoute
		if err := c.ShouldBindJSON(&woHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}

		workoutID, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		var workout shared.Workout

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(workoutID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with workout ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("workout")
		filter := bson.D{{Key: "_id", Value: id}}

		err = collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing user",
				"Exact": err.Error(),
			})
			return
		}

		if workout.UserID != userID {
			c.JSON(400, gin.H{
				"Error": "Issue with user in request",
				"Exact": errors.New("workout does not belong to provided user"),
			})
			return
		}

		if workout.Status != "Unstarted" {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout has already started, can't discard and retry"),
			})
			return
		}

		newworkout, err := WorkoutGen(woHandler.Time, woHandler.Difficulty, userID, database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}
		workoutRet := newworkout.(shared.Workout)

		patchUserGenerated(database, userID, true) // Even if it does error out, we don't want it to interrupt the WO generation process
		// The constantly running daily update code will check and update these anyways

		token := c.GetHeader("Authorization")

		resp, err := shared.PositionsRequestWorkout(workoutRet, token)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with positions API",
				"Exact": err.Error(),
			})
			return
		}

		if _, err := collection.DeleteOne(context.Background(), filter); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with deleting workout",
				"Exact": err,
			})
			return
		}

		c.JSON(201, gin.H{
			"workout":   workoutRet,
			"positions": resp,
		})

	}
}

func PostStretchWorkoutRetry(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.StrWorkoutRoute
		if err := c.ShouldBindJSON(&woHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}

		workoutID, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		var workout shared.StretchWorkout

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(workoutID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with workout ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("stretchworkout")
		filter := bson.D{{Key: "_id", Value: id}}

		err = collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing user",
				"Exact": err.Error(),
			})
			return
		}

		if workout.UserID != userID {
			c.JSON(400, gin.H{
				"Error": "Issue with user in request",
				"Exact": errors.New("workout does not belong to provided user"),
			})
			return
		}

		if workout.Status != "Unstarted" {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout has already started, can't discard and retry"),
			})
			return
		}

		newworkout, err := WorkoutGen(woHandler.Time, 0, userID, database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}
		workoutRet := newworkout.(shared.StretchWorkout)

		patchUserGenerated(database, userID, false) // Even if it does error out, we don't want it to interrupt the WO generation process
		// The constantly running daily update code will check and update these anyways

		token := c.GetHeader("Authorization")

		resp, err := shared.PositionsRequestStrWorkout(workoutRet, token)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with positions API",
				"Exact": err.Error(),
			})
			return
		}

		if _, err := collection.DeleteOne(context.Background(), filter); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with deleting workout",
				"Exact": err,
			})
			return
		}

		c.JSON(201, gin.H{
			"workout":   workoutRet,
			"positions": resp,
		})

	}
}

func patchUserGenerated(database *mongo.Database, userID string, isWO bool) error {
	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return err
	}

	collection := database.Collection("user")
	filter := bson.M{"_id": id}

	var update bson.M
	if isWO {
		update = bson.M{
			"$inc": bson.M{
				"wogenct":   1,
				"displevel": 2,
			},
		}
	} else {
		update = bson.M{
			"$inc": bson.M{
				"strwogenct": 1,
				"displevel":  2,
			},
		}
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
