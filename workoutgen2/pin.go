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

func PinStretchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.PinRoute
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
				"Error": "Issue with viewing str wo",
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

		pinCount, err := countPinnedWorkouts(database, userID, "stretchworkout")
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with counting pinned",
				"Exact": err.Error(),
			})
			return
		}

		if woHandler.Pinned && pinCount > 2 {
			c.JSON(400, gin.H{
				"Error": "Issue with pinning",
				"Exact": errors.New("should not be able to push a pin when more than 2 are pinened").Error(),
			})
			return
		}

		if workout.IsPinned != woHandler.Pinned {
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "pinned", Value: woHandler.Pinned},
				}},
			}

			_, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with update on db in request",
					"Exact": err.Error(),
				})
				return
			}
			c.Status(204)
		} else {
			c.JSON(200, gin.H{
				"Message": "Pinned",
			})
		}

	}
}

func PinWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var woHandler shared.PinRoute
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

		pinCount, err := countPinnedWorkouts(database, userID, "workout")
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with counting pinned",
				"Exact": err.Error(),
			})
			return
		}

		if woHandler.Pinned && pinCount > 2 {
			c.JSON(400, gin.H{
				"Error": "Issue with pinning",
				"Exact": errors.New("should not be able to push a pin when more than 2 are pinened").Error(),
			})
			return
		}

		if workout.IsPinned != woHandler.Pinned {
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "pinned", Value: woHandler.Pinned},
				}},
			}

			_, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with update on db in request",
					"Exact": err.Error(),
				})
				return
			}
			c.Status(204)
		} else {
			c.JSON(200, gin.H{
				"Message": "Pinned",
			})
		}
	}
}

func countPinnedWorkouts(database *mongo.Database, userID, collectionName string) (int64, error) {
	collection := database.Collection(collectionName)

	filterPinned := bson.D{
		{Key: "userid", Value: userID},
		{Key: "status", Value: bson.D{{Key: "$ne", Value: "Archived"}}},
		{Key: "pinned", Value: true},
	}

	count, err := collection.CountDocuments(context.Background(), filterPinned)
	if err != nil {
		return 0, err
	}

	return count, nil
}
