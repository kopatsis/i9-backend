package usergeneral

import (
	"context"
	"fmt"
	"fulli9/platform/middleware"
	"fulli9/ratings"
	"fulli9/shared"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		var exuser shared.User

		err = collection.FindOne(context.TODO(), filter).Decode(&exuser)
		if err == nil {
			refreshTokenDB(exuser.ID.Hex(), userBody.Token, database)
			c.Status(204)
			return
		}
		if err != mongo.ErrNoDocuments {
			c.JSON(400, gin.H{
				"Error": "Issue with checking user exists",
				"Exact": err.Error(),
			})
			return
		}

		user := shared.User{
			Username:          username,
			Name:              userBody.Name,
			Email:             userBody.Email,
			BirthMonth:        userBody.BirthMonth,
			BirthDay:          userBody.BirthDay,
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

		refreshTokenDB(result.InsertedID.(primitive.ObjectID).Hex(), userBody.Token, database)

		if userBody.Quizzed {
			ratings.PushQuiz(userBody.QuizRoute, database, result.InsertedID.(primitive.ObjectID).Hex())
		}

		// In both of the above cases, we do not want to error out for the actual creation if either of these processes themselves fail, so no err returned for them.

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

		var userBody shared.PatchUserRoute
		if err := c.ShouldBindJSON(&userBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		update := bson.M{
			"$set": bson.M{},
		}

		if userBody.Name != nil {
			update["$set"].(bson.M)["name"] = *userBody.Name
		}

		if userBody.Plyo != nil {
			update["$set"].(bson.M)["plyoToler"] = *userBody.Plyo
		}

		if userBody.Pushup != nil {
			update["$set"].(bson.M)["pushsetting"] = *userBody.Pushup
		}

		if userBody.BannedBody != nil {
			update["$set"].(bson.M)["bannedParts"] = *userBody.BannedBody
		}

		if userBody.Diff != nil {
			update["$set"].(bson.M)["lastdiff"] = *userBody.Diff
		}

		if userBody.Minutes != nil {
			update["$set"].(bson.M)["lastmins"] = *userBody.Minutes
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
			"$set": bson.M{
				"username": newAuthID,
				"name":     jwtBody.Name,
			},
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

		refreshTokenDB(mongoID, jwtBody.Refresh, database)
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
			if _, exists := c.GetQuery("force"); exists && err == mongo.ErrNoDocuments {
				createUserGet(database, c)
				return
			}
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

		user, err = userBadges(user, database)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with updating user badges",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(200, &user)

	}
}

func createUserGet(database *mongo.Database, c *gin.Context) {
	username, err := shared.GetSubFromJWT(c)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Issue with retrieving sub id",
			"Exact": err.Error(),
		})
		return
	}
	name := shared.GetNameIfExistsContext(c)

	user := shared.User{
		Username:          username,
		Name:              name,
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

	result, err := database.Collection("user").InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Issue adding user to database",
			"Exact": err.Error(),
		})
		return
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		c.JSON(400, gin.H{
			"Error": "Issue adding user to database",
			"Exact": "Created user id not primitive object id",
		})
		return
	}
	user.ID = oid

	c.JSON(200, &user)
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

func GetToken(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(200, gin.H{
				"idcode": "",
			})
			return
		}

		var dbToken shared.DBToken
		collection := database.Collection("usertoken")
		err = collection.FindOne(context.TODO(), bson.M{"user": userID}).Decode(&dbToken)
		if err != nil {
			c.JSON(200, gin.H{
				"idcode": "",
			})
			return
		}

		c.JSON(200, gin.H{
			"idcode": dbToken.ID.Hex(),
		})
	}
}

func refreshTokenDB(userid, refreshToken string, database *mongo.Database) error {
	collection := database.Collection("usertoken")

	filter := bson.M{"user": userid}
	update := bson.M{
		"$set": bson.M{
			"token": refreshToken,
		},
		"$setOnInsert": bson.M{
			"_id":  primitive.NewObjectID(),
			"user": userid,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	if err != nil {
		return fmt.Errorf("failed to update or insert token: %v", err)
	}

	return nil
}

func userBadges(user shared.User, database *mongo.Database) (shared.User, error) {

	changed := false
	badges := user.Badges

	levelBdg := "Achieved a Progress Level of " + roundDownLevel(user.DisplayLevel)
	if !slices.Contains(badges, levelBdg) {
		changed = true
		badges = replaceInUser("Achieved a Progress Level of ", levelBdg, badges)
	}

	woGenBdg := "Workouts: Generated at least " + roundDownCount(user.WOGeneratedCt)
	if !slices.Contains(badges, woGenBdg) {
		changed = true
		badges = replaceInUser("Workouts: Generated at least ", woGenBdg, badges)
	}

	stwoGenBdg := "Stretch Workouts: Generated at least " + roundDownCount(user.StrWOGeneratedCt)
	if !slices.Contains(badges, stwoGenBdg) {
		changed = true
		badges = replaceInUser("Stretch Workouts: Generated at least ", stwoGenBdg, badges)
	}

	ratedBdg := "Workouts: Rated at least " + roundDownCount(user.WORatedCt)
	if !slices.Contains(badges, ratedBdg) {
		changed = true
		badges = replaceInUser("Workouts: Rated at least ", ratedBdg, badges)
	}

	woStBdg := "Workouts: Started at least " + roundDownCount(user.WOStartedCt)
	if !slices.Contains(user.Badges, woStBdg) {
		changed = true
		badges = replaceInUser("Workouts: Started at least ", woStBdg, badges)
	}

	stwoStBdg := "Stretch Workouts: Started at least " + roundDownCount(user.StrWOGeneratedCt)
	if !slices.Contains(badges, stwoStBdg) {
		changed = true
		badges = replaceInUser("Stretch Workouts: Started at least ", stwoStBdg, badges)
	}

	if !changed {
		return user, nil
	} else {
		collection := database.Collection("user")
		filter := bson.M{"_id": user.ID}
		update := bson.M{
			"$set": bson.M{
				"badges": badges,
			},
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updatedUser shared.User
		err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedUser)
		if err != nil {
			return updatedUser, err
		}

		return updatedUser, nil
	}

}

func replaceInUser(prefix, phrase string, slice []string) []string {
	new := []string{}

	for _, st := range slice {
		if !strings.Contains(st, prefix) {
			new = append(new, st)
		}
	}

	if phrase != prefix {
		new = append(new, phrase)
	}

	return new
}

func roundDownCount(ct int) string {
	thresholds := []int{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000}
	rounded := 1 // Default to the smallest value

	if ct < rounded {
		return ""
	}

	for _, threshold := range thresholds {
		if ct < threshold {
			break
		}
		rounded = threshold
	}

	return strconv.Itoa(rounded)
}

func roundDownLevel(ct int) string {
	thresholds := []int{50, 100, 250, 500, 750, 1000, 1250, 1500, 2000, 2500, 3000, 3500, 4000, 5000, 6000, 7000, 8000, 10000, 12500, 15000, 20000, 25000, 30000, 40000, 50000, 75000, 100000}
	rounded := 50 // Default to the smallest value in the list

	if ct < rounded {
		return ""
	}

	for _, threshold := range thresholds {
		if ct < threshold {
			break
		}
		rounded = threshold
	}

	return strconv.Itoa(rounded)
}
