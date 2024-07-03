package views

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

func GetStrByID(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

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

		strs, err := shared.GetStretchHelper(database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": err.Error(),
			})
			return
		}

		for _, str := range strs {
			if str.ID == id {
				c.JSON(200, &str)
				return
			}
		}

		c.JSON(400, gin.H{
			"Error": "Issue with viewing stretch",
			"Exact": errors.New("does not exist with provided id"),
		})
	}
}

func GetStrecthes(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		strtype := c.DefaultQuery("type", "")

		var filterStr bson.M

		if strtype != "" {
			filterStr = bson.M{"status": strtype}
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
