package adapts

import (
	"context"
	"fulli9/shared"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CloneWorkoutHandler(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var workout shared.Workout

		idStr, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL parameter",
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("workout")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workout",
				"Exact": err.Error(),
			})
			return
		}

		user, err := getUserHelper(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing user",
				"Exact": err.Error(),
			})
			return
		}

		workout.Date = primitive.NewDateTimeFromTime(time.Now())
		workout.ID = primitive.NilObjectID
		workout.UserID = user.ID.Hex()
		workout.Username = user.Username
		workout.Status = "Unstarted"
		workout.LevelAtStart = user.Level

		collection = database.Collection("workout")
		result, err := collection.InsertOne(context.Background(), workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with posting new workout",
				"Exact": err.Error(),
			})
			return
		}
		workout.ID = result.InsertedID.(primitive.ObjectID)

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

func CloneStretchWorkoutHandler(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var workout shared.StretchWorkout

		idStr, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL parameter",
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("stretchworkout")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workout",
				"Exact": err.Error(),
			})
			return
		}

		user, err := getUserHelper(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing user",
				"Exact": err.Error(),
			})
			return
		}

		workout.Date = primitive.NewDateTimeFromTime(time.Now())
		workout.ID = primitive.NilObjectID
		workout.UserID = user.ID.Hex()
		workout.Status = "Unstarted"
		workout.LevelAtStart = user.Level

		collection = database.Collection("stretchworkout")
		result, err := collection.InsertOne(context.Background(), workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with posting new workout",
				"Exact": err.Error(),
			})
			return
		}
		workout.ID = result.InsertedID.(primitive.ObjectID)

		token := c.GetHeader("Authorization")

		resp, err := shared.PositionsRequestStrWorkout(workout, token)
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

func getUserHelper(database *mongo.Database, c *gin.Context) (shared.User, error) {
	var user shared.User

	userID, err := shared.GetIDFromReq(database, c)
	if err != nil {
		return user, err
	}

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return user, err
	}

	collection := database.Collection("user")
	filter := bson.D{{Key: "_id", Value: id}}

	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
