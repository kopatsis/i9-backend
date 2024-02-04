package userfuncs

import (
	"context"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ClearBannedExer(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.ExerListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userBody.UserID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		newBannedExers := []string{}

		update := bson.M{
			"$set": bson.M{"bannedExer": newBannedExers},
		}

		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                   userBody.UserID,
			"Banned Exercise List": newBannedExers,
		})
	}
}

func ClearBannedStr(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.StrListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userBody.UserID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		newBannedStrs := []string{}

		update := bson.M{
			"$set": bson.M{"bannedStr": newBannedStrs},
		}

		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                  userBody.UserID,
			"Banned Strecth List": newBannedStrs,
		})
	}
}

func ClearBannedBody(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.BodyListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userBody.UserID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		newBannedBodyParts := []int{}

		update := bson.M{
			"$set": bson.M{"bannedParts": newBannedBodyParts},
		}

		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                     userBody.UserID,
			"Banned Body Parts List": newBannedBodyParts,
		})
	}
}

func ClearExerFav(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.ExerListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userBody.UserID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		newExerFavorites := map[string]float32{}

		update := bson.M{
			"$set": bson.M{"exerfavs": newExerFavorites},
		}

		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                             userBody.UserID,
			"Banned Exercise Favorite Rates": newExerFavorites,
		})
	}
}
