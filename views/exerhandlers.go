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

func GetExerByID(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
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
				"Error": "Issue with exercise ID",
				"Exact": err.Error(),
			})
			return
		}

		exers, err := shared.GetExersHelper(database, boltDB)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exer",
				"Exact": err.Error(),
			})
			return
		}

		for _, exer := range exers {
			if exer.ID == id {
				c.JSON(200, &exer)
				return
			}
		}

		c.JSON(400, gin.H{
			"Error": "Issue with viewing exer",
			"Exact": errors.New("does not exist with provided id").Error(),
		})

	}
}

func GetExercises(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		exertype := c.DefaultQuery("type", "")

		var filterEx bson.M

		if exertype != "" {
			filterEx = bson.M{"parent": exertype}
		}

		collection := database.Collection("exercise")

		cursor, err := collection.Find(context.Background(), filterEx)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exercises",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastExers []shared.Exercise
		err = cursor.All(context.Background(), &pastExers)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exercises",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastExers) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exercises",
				"Exact": "No results returned (check type)",
			})
			return
		}

		c.JSON(200, &pastExers)

	}
}
