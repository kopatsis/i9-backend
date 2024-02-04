package views

import (
	"context"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStrByID(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var stretch shared.Stretch

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
				"Error": "Issue with stretch ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("stretch")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&stretch)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, &stretch)

	}
}

func GetStrecthes(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		strtype := c.DefaultQuery("type", "")

		var filterStr primitive.D

		if strtype != "" {
			filterStr = bson.D{
				{Key: "status", Value: strtype},
			}
		}

		collection := database.Collection("stretch")

		cursor, err := collection.Find(context.Background(), filterStr)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastStrs []shared.Stretch
		err = cursor.All(context.Background(), &pastStrs)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastStrs) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": "No results returned (check type)",
			})
			return
		}

		c.JSON(200, &pastStrs)

	}
}
