package adapts

import (
	"errors"
	"fulli9/adapts/alteredfuncs"
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostRestartedWorkout(database *mongo.Database, boltDB *bbolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		workout, err := Adapt(-1, userID, database, boltDB, id, false, false)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with adapted workout generator",
				"Exact": err.Error(),
			})
			return
		}

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workout)
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

func PostRestartedStrWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		workout, err := alteredfuncs.GetStretchWorkoutByID(database, id)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with fetching wo",
				"Exact": err.Error(),
			})
			return
		}

		if workout.UserID != userID {
			c.JSON(400, gin.H{
				"Error": "Issue with matching wo",
				"Exact": errors.New("workout does not match the user that's requesting it"),
			})
			return
		}

		workout.Status = "Unstarted"
		workout.PausedTime = 0

		if err := alteredfuncs.SaveUpdatedStretchWorkout(database, workout); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue saving wo",
				"Exact": err.Error(),
			})
			return
		}

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workout)
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

func PostRestartedIntroWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		workout, err := alteredfuncs.GetWorkoutByID(database, id)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with fetching wo",
				"Exact": err.Error(),
			})
			return
		}

		if workout.UserID != userID {
			c.JSON(400, gin.H{
				"Error": "Issue with matching wo",
				"Exact": errors.New("workout does not match the user that's requesting it"),
			})
			return
		}

		workout.Status = "Unstarted"
		workout.PausedTime = 0

		if err := alteredfuncs.SaveUpdatedWorkout(database, workout); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue saving wo",
				"Exact": err.Error(),
			})
			return
		}

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workout)
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
