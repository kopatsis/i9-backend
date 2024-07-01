package views

import (
	"context"
	"fulli9/shared"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Will split to auth and unauth once login developed

		var workout shared.Workout

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
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("workout")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workout",
				"Exact": err.Error(),
			})
			return
		}

		if _, exists := c.GetQuery("script"); !exists {
			c.JSON(200, &workout)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestWorkout(workout, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workout,
				"positions": resp,
			})
		}
	}
}

func GetWorkouts(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		limit := c.DefaultQuery("limit", "")

		optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

		var limitNum int64
		if limit != "" {
			var err error
			limitNum, err = strconv.ParseInt(limit, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter limit conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetLimit(limitNum)
		}

		skips := c.DefaultQuery("skips", "")

		var skipNum int64
		if skips != "" {
			var err error
			skipNum, err = strconv.ParseInt(skips, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter skip conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetSkip(skipNum)
		}

		collection := database.Collection("workout")

		filterWO := bson.D{
			{Key: "userid", Value: userID},
		}

		if _, exists := c.GetQuery("includearch"); !exists {
			filterWO = append(filterWO, bson.E{Key: "status", Value: bson.D{{Key: "$ne", Value: "Archived"}}})
		}

		cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastWorkouts []shared.Workout
		err = cursor.All(context.Background(), &pastWorkouts)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastWorkouts) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": "No results returned",
			})
			return
		}

		if int64(len(pastWorkouts)) < limitNum || limitNum == 0 {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  false,
				"Next Skips": 0,
				"Result":     &pastWorkouts,
			})
		} else {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  true,
				"Next Skips": skipNum + limitNum,
				"Result":     &pastWorkouts,
			})
		}

	}
}

func GetStretchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Will split to auth and unauth once login developed

		var workout shared.StretchWorkout

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
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("stretchworkout")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workout",
				"Exact": err.Error(),
			})
			return
		}

		if _, exists := c.GetQuery("script"); !exists {
			c.JSON(200, &workout)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestStrWorkout(workout, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workout,
				"positions": resp,
			})
		}

	}
}

func GetStretchWorkouts(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		limit := c.DefaultQuery("limit", "")

		optionsWO := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

		var limitNum int64
		if limit != "" {
			var err error
			limitNum, err = strconv.ParseInt(limit, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter limit conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetLimit(limitNum)
		}

		skips := c.DefaultQuery("skips", "")

		var skipNum int64
		if skips != "" {
			var err error
			skipNum, err = strconv.ParseInt(skips, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with query parameter skip conversion",
					"Exact": err.Error(),
				})
				return
			}
			optionsWO = optionsWO.SetSkip(skipNum)
		}

		collection := database.Collection("stretchworkout")

		filterWO := bson.D{
			{Key: "userid", Value: userID},
		}

		if _, exists := c.GetQuery("includearch"); !exists {
			filterWO = append(filterWO, bson.E{Key: "status", Value: bson.D{{Key: "$ne", Value: "Archived"}}})
		}

		cursor, err := collection.Find(context.Background(), filterWO, optionsWO)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastWorkouts []shared.StretchWorkout
		err = cursor.All(context.Background(), &pastWorkouts)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastWorkouts) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": "No results returned",
			})
			return
		}

		if int64(len(pastWorkouts)) < limitNum || limitNum == 0 {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  false,
				"Next Skips": 0,
				"Result":     &pastWorkouts,
			})
		} else {
			c.JSON(200, gin.H{
				"Found":      len(pastWorkouts),
				"Remaining":  true,
				"Next Skips": skipNum + limitNum,
				"Result":     &pastWorkouts,
			})
		}

	}
}

func GetMostRecent(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		optionsWO := options.FindOne().SetSort(bson.D{{Key: "date", Value: -1}})

		filterWO := bson.D{
			{Key: "userid", Value: userID},
			{Key: "status", Value: bson.D{{Key: "$ne", Value: "Archived"}}},
		}

		var st shared.StretchWorkout
		var errSt error
		var wo shared.Workout
		var errWo error

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			collection := database.Collection("stretchworkout")
			errSt = collection.FindOne(context.Background(), filterWO, optionsWO).Decode(&st)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			collection := database.Collection("workout")
			errWo = collection.FindOne(context.Background(), filterWO, optionsWO).Decode(&wo)
		}()

		wg.Wait()

		if errSt != nil && errSt != mongo.ErrNoDocuments {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch workouts",
				"Exact": errSt.Error(),
			})
			return
		}

		if errWo != nil && errWo != mongo.ErrNoDocuments {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workouts",
				"Exact": errWo.Error(),
			})
			return
		}

		if errWo == mongo.ErrNoDocuments && errSt == mongo.ErrNoDocuments {
			c.Status(204)
			return
		}

		if errWo == mongo.ErrNoDocuments || wo.Created.Time().Before(st.Created.Time()) {
			c.JSON(200, gin.H{
				"name":   st.Name,
				"id":     st.ID,
				"date":   st.Created,
				"status": st.Status,
				"type":   "Stretch",
				"stored": false,
			})
			return
		}

		woType := "Regular"
		if wo.IsIntro {
			woType = "Intro"
		}

		c.JSON(200, gin.H{
			"name":   wo.Name,
			"id":     wo.ID,
			"date":   wo.Created,
			"status": wo.Status,
			"type":   woType,
			"stored": false,
		})
	}
}
