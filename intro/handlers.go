package intro

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

func PostIntroWorkout(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.IntroWorkoutRoute
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
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		workout, err := GenerateIntroWorkout(woHandler.Time, userID, database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workout)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestWorkout(workout, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workout,
				"positions": resp,
			})
		}
	}
}

func PostIntroWorkoutRetry(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.IntroWorkoutRoute
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
				"Exact": "Unable to get ID from URL parameter",
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
				"Error": "Issue with viewing wo",
				"Exact": err.Error(),
			})
			return
		}

		if workout.UserID != userID {
			c.JSON(400, gin.H{
				"Error": "Issue with user in request",
				"Exact": errors.New("workout does not belong to provided user").Error(),
			})
			return
		}

		if workout.Status != "Unstarted" {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout has already started, can't discard and retry").Error(),
			})
			return
		}

		workoutRet, err := GenerateIntroWorkout(woHandler.Time, userID, database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}

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
