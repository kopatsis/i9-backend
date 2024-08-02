package workoutgen2

import (
	"context"
	"errors"
	"fulli9/shared"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PatchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.PatchWorkout
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

		if workout.Minutes < woHandler.PausedMinutes {
			woHandler.PausedMinutes = workout.Minutes
		}

		allowedStatus := []string{"Unstarted", "Progressing", "Paused", "Completed", "Rated", "Archived"}

		if slices.Contains(allowedStatus, workout.Status) {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)").Error(),
			})
			return
		}

		if (woHandler.Status == "Progressing" || woHandler.Status == "Paused") && workout.Status == "Unstarted" {
			workout.StartedCount++
			workout.LastStarted = primitive.NewDateTimeFromTime(time.Now())
			workout.StartedDates = append(workout.StartedDates, workout.LastStarted)
			if err := patchUserStarted(database, userID, true); err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with updating user ct",
					"Exact": err.Error(),
				})
				return
			}
		}

		if woHandler.Status == "Completed" && workout.Status != "Completed" {
			workout.FinishedCount++
			if err := patchUserFinished(database, userID, true); err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with updating user ct",
					"Exact": err.Error(),
				})
				return
			}
		}

		workout.PausedTime = woHandler.PausedMinutes
		workout.Status = woHandler.Status

		updateFilter := bson.M{"_id": workout.ID}

		update := bson.M{"$set": workout}

		_, err = collection.UpdateOne(context.Background(), updateFilter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving update",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)").Error(),
			})
			return
		}

		c.Status(204)
	}
}

func PatchStretchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.PatchWorkout
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

		if workout.Minutes < woHandler.PausedMinutes {
			woHandler.PausedMinutes = workout.Minutes
		}

		allowedStatus := []string{"Unstarted", "Progressing", "Paused", "Completed", "Rated", "Archived"}

		if slices.Contains(allowedStatus, workout.Status) {
			c.JSON(400, gin.H{
				"Error": "Issue with status in request",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)").Error(),
			})
			return
		}

		if (woHandler.Status == "Progressing" || woHandler.Status == "Paused") && workout.Status == "Unstarted" {
			workout.StartedCount++
			workout.LastStarted = primitive.NewDateTimeFromTime(time.Now())
			workout.StartedDates = append(workout.StartedDates, workout.LastStarted)
			if err := patchUserStarted(database, userID, false); err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with updating user ct",
					"Exact": err.Error(),
				})
				return
			}
		}

		if woHandler.Status == "Completed" && workout.Status != "Completed" {
			workout.FinishedCount++
			if err := patchUserFinished(database, userID, false); err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with updating user ct",
					"Exact": err.Error(),
				})
				return
			}
		}

		workout.PausedTime = woHandler.PausedMinutes
		workout.Status = woHandler.Status

		updateFilter := bson.M{"_id": workout.ID}

		update := bson.M{"$set": workout}

		_, err = collection.UpdateOne(context.Background(), updateFilter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving update",
				"Exact": errors.New("workout status not in allowable values (Unstarted, Paused, Rated, Archived)").Error(),
			})
			return
		}

		c.Status(204)
	}
}

func Rename(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request shared.RenameRoute

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		if request.Type != "exercise" && request.Type != "stretch" {
			c.JSON(400, gin.H{
				"Error": "Wrong type",
				"Exact": "Wrong type",
			})
			return
		}

		if len(request.Name) > 100 {
			request.Name = request.Name[0:99]
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

		if request.Type == "stretch" {
			collection = database.Collection("stretchworkout")
			var strwo shared.StretchWorkout

			err = collection.FindOne(context.Background(), filter).Decode(&strwo)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with viewing wo",
					"Exact": err.Error(),
				})
				return
			}

			if strwo.UserID != userID {
				c.JSON(400, gin.H{
					"Error": "Issue with user in request",
					"Exact": errors.New("workout does not belong to provided user").Error(),
				})
				return
			}

		} else {
			var workout shared.Workout

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

		}

		update := bson.M{"$set": bson.M{"name": request.Name}}

		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving update",
				"Exact": err,
			})
			return
		}

		c.Status(204)
	}
}

func patchUserStarted(database *mongo.Database, userID string, isWO bool) error {
	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return err
	}

	collection := database.Collection("user")
	filter := bson.M{"_id": id}

	var user shared.User
	if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return err
	}

	var update string
	if isWO {
		update = "wostartct"
	} else {
		update = "strwostartct"
	}

	if err := IncrementMonthly(user, database, update); err != nil {
		return err
	}

	if err := IncrementDispLevelBy(user, database, 2); err != nil {
		return err
	}

	return nil
}

func patchUserFinished(database *mongo.Database, userID string, isWO bool) error {
	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return err
	}

	collection := database.Collection("user")
	filter := bson.M{"_id": id}

	var user shared.User
	if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return err
	}

	var update string
	if isWO {
		update = "wofinishct"
	} else {
		update = "strwofinishct"
	}

	if err := IncrementMonthly(user, database, update); err != nil {
		return err
	}

	if err := IncrementDispLevelBy(user, database, 2); err != nil {
		return err
	}

	return nil
}

func IncrementDispLevelBy(user shared.User, database *mongo.Database, increment int) error {

	user.DisplayLevel += increment

	loc, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(loc)
	current := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var lastTime time.Time
	if len(user.LevelHistory) > 0 {
		lastTime = user.LevelHistory[len(user.LevelHistory)-1].Date.Time()
	}

	if current.Equal(lastTime) {
		user.LevelHistory[len(user.LevelHistory)-1].Level = user.DisplayLevel
	} else {
		user.LevelHistory = append(user.LevelHistory, shared.LevelHistory{
			Date:  primitive.NewDateTimeFromTime(current),
			Level: user.DisplayLevel,
		})
	}

	collection := database.Collection("user")
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func IncrementMonthly(user shared.User, database *mongo.Database, field string) error {

	loc, _ := time.LoadLocation("America/Los_Angeles")
	monthInt := int(time.Now().In(loc).Month()) - 1
	year := time.Now().In(loc).Year()

	if user.MonthlyHistory[monthInt].Year != year {
		user.MonthlyHistory[monthInt] = shared.MonthlyHistory{
			Year: year,
		}
	}

	switch field {
	case "strwostartct":
		user.MonthlyHistory[monthInt].StrWOStartedCt++
		user.StrWOStartedCt++
	case "wogenct":
		user.MonthlyHistory[monthInt].WOGeneratedCt++
		user.StrWOStartedCt++
	case "strwogenct":
		user.MonthlyHistory[monthInt].StrWOGeneratedCt++
		user.WOGeneratedCt++
	case "wostartct":
		user.MonthlyHistory[monthInt].WOStartedCt++
		user.WOStartedCt++
	case "completed":
		user.MonthlyHistory[monthInt].WORatedCt++
		user.WORatedCt++
	case "strwocompleted":
		user.MonthlyHistory[monthInt].StrWORatedCt++
		user.StrWORatedCt++
	case "wofinishct":
		user.MonthlyHistory[monthInt].WOFinishedCt++
		user.WOFinishedCt++
	case "strwofinishct":
		user.MonthlyHistory[monthInt].StrWOFinishedCt++
		user.StrWOFinishedCt++
	}

	collection := database.Collection("user")
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
