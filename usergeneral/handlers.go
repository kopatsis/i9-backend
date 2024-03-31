package usergeneral

import (
	"context"
	"fulli9/platform/middleware"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.UserRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		username, err := shared.GetSubFromJWT(c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with retrieving sub id",
				"Exact": err.Error(),
			})
			return
		}

		filter := bson.M{"username": username}

		collection := database.Collection("user")

		if err := collection.FindOne(context.TODO(), filter).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(204)
				return
			} else {
				c.JSON(400, gin.H{
					"Error": "Issue with checking user exists",
					"Exact": err.Error(),
				})
				return
			}
		}

		user := shared.User{
			Username:          username,
			Name:              userBody.Name,
			PlyoTolerance:     3,
			PushupSetting:     "Knee",
			BannedExercises:   []string{},
			BannedStretches:   []string{},
			BannedParts:       []int{},
			ExerFavoriteRates: map[string]float32{},
			ExerModifications: map[string]float32{},
			TypeModifications: map[string]float32{},
			RoundEndurance:    map[int]float32{},
			TimeEndurance:     map[int]float32{},
		}

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

func PostLocalUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.UserRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		user := shared.User{
			Username:          "",
			Name:              userBody.Name,
			PlyoTolerance:     3,
			BannedExercises:   []string{},
			BannedStretches:   []string{},
			BannedParts:       []int{},
			PushupSetting:     "Knee",
			ExerFavoriteRates: map[string]float32{},
			ExerModifications: map[string]float32{},
			TypeModifications: map[string]float32{},
			RoundEndurance:    map[int]float32{},
			TimeEndurance:     map[int]float32{},
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

		update := bson.M{
			"$set": bson.M{"username": result.InsertedID.(primitive.ObjectID).Hex()},
		}

		filter := bson.M{"_id": result.InsertedID}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating newly created local user",
				"Exact": err.Error(),
			})
			return
		}

		tokenString, err := middleware.GenerateJWT(result.InsertedID.(primitive.ObjectID).Hex())
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with local user JWT gen",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"ID":    result.InsertedID,
			"Token": tokenString,
		})
	}
}

func GetLocalJWT(database *mongo.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr, exists := ctx.Params.Get("id")
		if !exists {
			ctx.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL parameter",
			})
			return
		}

		userID, err := shared.AuthIDtoMongoID(database, idStr)
		if err != nil {
			ctx.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		if idStr != userID {
			ctx.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Non-Local user",
			})
			return
		}

		tokenString, err := middleware.GenerateJWT(idStr)
		if err != nil {
			ctx.JSON(400, gin.H{
				"Error": "Issue with local user JWT retrieval",
				"Exact": err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"JWT": tokenString,
		})
	}
}

func PatchUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.UserRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		if userBody.Name == "" {
			c.Status(204)
			return
		}

		update := bson.M{
			"$set": bson.M{"name": userBody.Name},
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		filter := bson.M{"_id": id}

		collection := database.Collection("user")
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID": userID,
		})
	}
}

func MergeLocalUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var jwtBody shared.MergeRoute
		if err := c.ShouldBindJSON(&jwtBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		sub, err := shared.GetSubFromJWTStr(jwtBody.LocalJWT)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with old JWT id",
				"Exact": err.Error(),
			})
			return
		}

		mongoID, err := shared.AuthIDtoMongoID(database, sub)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with old JWT id to mongo id",
				"Exact": err.Error(),
			})
			return
		}

		newAuthID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with new JWT id",
				"Exact": err.Error(),
			})
			return
		}

		update := bson.M{
			"$set": bson.M{"username": newAuthID},
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(mongoID); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		filter := bson.M{"_id": id}

		collection := database.Collection("user")
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user to new id",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID": mongoID,
		})
	}
}

func GetUser(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user shared.User

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
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

		err = collection.FindOne(context.Background(), filter).Decode(&user)
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

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
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
