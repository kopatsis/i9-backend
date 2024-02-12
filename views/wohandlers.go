package views

import (
	"context"
	"fulli9/shared"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Will split to auth and unauth once login developed

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

		c.JSON(200, &workout)

	}
}

func GetWorkouts(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userHandler shared.UserIDRoute
		if err := c.ShouldBindJSON(&userHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		limit := c.DefaultQuery("limit", "")

		optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

		var limitNum int64
		if limit != "" {
			var err error
			limitNum, err = strconv.ParseInt(limit, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter limit conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetLimit(limitNum)
		}

		skips := c.DefaultQuery("skips", "")

		var skipNum int64
		if skips != "" {
			var err error
			skipNum, err = strconv.ParseInt(skips, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter skip conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetSkip(skipNum)
		}

		collection := database.Collection("workout")

		filterWO := bson.D{
			{Key: "userid", Value: userHandler.UserID},
		}

		cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastWorkouts []shared.Workout
		err = cursor.All(context.Background(), &pastWorkouts)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastWorkouts) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": "No results returned",
			})
			return
		}

		if int64(len(pastWorkouts)) < limitNum || limitNum == 0 {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  false,
				"Next Skips": 0,
				"Result":     &pastWorkouts,
			})
		} else {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  true,
				"Next Skips": skipNum + limitNum,
				"Result":     &pastWorkouts,
			})
		}

	}
}

func GetStretchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Will split to auth and unauth once login developed

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

		collection := database.Collection("stretchworkout")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workout",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, &workout)

	}
}

func GetStretchWorkouts(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userHandler shared.UserIDRoute
		if err := c.ShouldBindJSON(&userHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		limit := c.DefaultQuery("limit", "")

		optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

		var limitNum int64
		if limit != "" {
			var err error
			limitNum, err = strconv.ParseInt(limit, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter limit conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetLimit(limitNum)
		}

		skips := c.DefaultQuery("skips", "")

		var skipNum int64
		if skips != "" {
			var err error
			skipNum, err = strconv.ParseInt(skips, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter skip conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetSkip(skipNum)
		}

		collection := database.Collection("stretchworkout")

		filterWO := bson.D{
			{Key: "userid", Value: userHandler.UserID},
		}

		cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastWorkouts []shared.Workout
		err = cursor.All(context.Background(), &pastWorkouts)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastWorkouts) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": "No results returned",
			})
			return
		}

		if int64(len(pastWorkouts)) < limitNum || limitNum == 0 {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  false,
				"Next Skips": 0,
				"Result":     &pastWorkouts,
			})
		} else {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  true,
				"Next Skips": skipNum + limitNum,
				"Result":     &pastWorkouts,
			})
		}

	}
}
