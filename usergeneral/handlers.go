package usergeneral

import (
	"context"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.PostUserRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		user := shared.User{
			Username: userBody.UserName,
			Name:     userBody.Name,
		}

		collection := database.Collection("user")
		result, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue adding user to database",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"ID": result.InsertedID,
		})
	}
}

func PutUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.PutUserRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		if userBody.Name == "" && userBody.UserName == "" {
			c.Status(204)
			return
		}

		var update primitive.M
		if userBody.Name != "" {
			update = bson.M{
				"$set": bson.M{"username": userBody.UserName},
			}
		} else if userBody.UserName == "" {
			update = bson.M{
				"$set": bson.M{"name": userBody.Name},
			}
		} else {
			update = bson.M{
				"$set": bson.M{"name": userBody.Name, "username": userBody.UserName},
			}
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

		filter := bson.M{"id": id}

		collection := database.Collection("user")
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID": userBody.UserID,
		})
	}
}

func GetUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user shared.User

		var userBody shared.UserIDRoute
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

		err := collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, &user)

	}
}

func DeleteUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.UserIDRoute
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

		result, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with deleting user",
				"Exact": err.Error(),
			})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with deleting user",
				"Exact": "No user was deleted",
			})
			return
		}

		c.Status(204)

	}
}
