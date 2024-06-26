package workoutgen2

import (
	"context"
	"errors"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PatchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.PatchWorkout
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

		if workout.Minutes < woHandler.PausedMinutes {
			woHandler.PausedMinutes = workout.Minutes
		}

		if woHandler.Status != "Unstarted" && woHandler.Status != "Progressing" && woHandler.Status != "Paused" && woHandler.Status != "Rated" && woHandler.Status != "Archived" {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)"),
			})
			return
		}

		workout.PausedTime = woHandler.PausedMinutes
		workout.Status = woHandler.Status

		updateFilter := bson.M{"_id": workout.ID}

		update := bson.M{"$set": workout}

		_, err = collection.UpdateOne(context.Background(), updateFilter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving update",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)"),
			})
			return
		}

		c.Status(204)
	}
}

func PatchStretchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.PatchWorkout
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

		if workout.Minutes < woHandler.PausedMinutes {
			woHandler.PausedMinutes = workout.Minutes
		}

		if woHandler.Status != "Unstarted" && woHandler.Status != "Progressing" && woHandler.Status != "Paused" && woHandler.Status != "Finished" && woHandler.Status != "Archived" {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)"),
			})
			return
		}

		workout.PausedTime = woHandler.PausedMinutes
		workout.Status = woHandler.Status

		updateFilter := bson.M{"_id": workout.ID}

		update := bson.M{"$set": workout}

		_, err = collection.UpdateOne(context.Background(), updateFilter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving update",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)"),
			})
			return
		}

		c.Status(204)
	}
}

func Rename(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request shared.RenameRoute

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		if request.Type != "exercise" && request.Type != "stretch" {
			c.JSON(400, gin.H{
				"Error": "Wrong type",
				"Exact": "Wrong type",
			})
			return
		}

		if len(request.Name) > 100 {
			request.Name = request.Name[0:99]
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

		if request.Type == "stretch" {
			collection = database.Collection("stretchworkout")
			var strwo shared.StretchWorkout

			err = collection.FindOne(context.Background(), filter).Decode(&strwo)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with viewing wo",
					"Exact": err.Error(),
				})
				return
			}

			if strwo.UserID != userID {
				c.JSON(400, gin.H{
					"Error": "Issue with user in request",
					"Exact": errors.New("workout does not belong to provided user"),
				})
				return
			}

		} else {
			var workout shared.Workout

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

		}

		update := bson.M{"$set": bson.M{"name": request.Name}}

		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving update",
				"Exact": err,
			})
			return
		}

		c.Status(204)

	}
}
