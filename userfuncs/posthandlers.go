package userfuncs

import (
	"context"
	"fulli9/shared"
	"slices"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func uniqueStrList(list1, list2 []string) []string {
	new := list1
	for _, val := range list2 {
		if !slices.Contains(new, val) {
			new = append(new, val)
		}
	}
	return new
}

func uniqueIntList(list1, list2 []int) []int {
	new := list1
	for _, val := range list2 {
		if !slices.Contains(new, val) {
			new = append(new, val)
		}
	}
	return new
}

func PostPushupSetting(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.PushupSettingRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		if userBody.Setting != "Wall" && userBody.Setting != "Knee" && userBody.Setting != "Regular" {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": "provided setting for pushups non-existentx",
			})
			return
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		update := bson.M{
			"$set": bson.M{"pushsetting": userBody.Setting},
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

func PostPlyo(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userBody shared.PlyoRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		update := bson.M{
			"$set": bson.M{"plyoToler": userBody.Plyo},
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

func PostBannedExer(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user shared.User

		var userBody shared.ExerListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
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

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with finding user",
				"Exact": err.Error(),
			})
			return
		}

		newBannedExers := uniqueStrList(user.BannedExercises, userBody.ExerList)

		update := bson.M{
			"$set": bson.M{"bannedExer": newBannedExers},
		}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                   userID,
			"Banned Exercise List": newBannedExers,
		})
	}
}

func PostBannedStr(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user shared.User

		var userBody shared.StrListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
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

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with finding user",
				"Exact": err.Error(),
			})
			return
		}

		newBannedStrs := uniqueStrList(user.BannedStretches, userBody.StrList)

		update := bson.M{
			"$set": bson.M{"bannedStr": newBannedStrs},
		}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                  userID,
			"Banned Stretch List": newBannedStrs,
		})
	}
}

func PostBannedBody(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user shared.User

		var userBody shared.BodyListRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
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

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with finding user",
				"Exact": err.Error(),
			})
			return
		}

		newBannedBodyParts := uniqueIntList(user.BannedParts, userBody.BodyList)

		update := bson.M{
			"$set": bson.M{"bannedParts": newBannedBodyParts},
		}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                     userID,
			"Banned Body Parts List": newBannedBodyParts,
		})
	}
}

func PostExerFav(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user shared.User

		var userBody shared.ExerMapRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
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

		collection := database.Collection("user")
		filter := bson.D{{Key: "_id", Value: id}}

		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with finding user",
				"Exact": err.Error(),
			})
			return
		}

		newExerFavorites := map[string]float32{}

		if user.ExerFavoriteRates != nil {
			newExerFavorites = user.ExerFavoriteRates
		}

		for key, val := range userBody.ExerMap {
			newExerFavorites[key] = val
		}

		update := bson.M{
			"$set": bson.M{"exerfavs": newExerFavorites},
		}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"ID":                             userID,
			"Banned Exercise Favorite Rates": newExerFavorites,
		})
	}
}
