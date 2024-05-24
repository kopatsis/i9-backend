package workoutgen2

import (
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.WorkoutRoute
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

		workout, err := WorkoutGen(woHandler.Time, woHandler.Difficulty, userID, database)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}
		workoutRet := workout.(shared.Workout)

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workoutRet)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestWorkout(workoutRet, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workoutRet,
				"positions": resp,
			})
		}

	}
}

func PostStretchWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.StrWorkoutRoute
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
				"Error": "Issue with stretch workout generator",
				"Exact": err.Error(),
			})
			return
		}

		workout, err := WorkoutGen(woHandler.Time, 0, userID, database)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with stretch workout generator",
				"Exact": err.Error(),
			})
			return
		}
		workoutRet := workout.(shared.StretchWorkout)

		if _, exists := c.GetQuery("noscript"); exists {
			c.JSON(201, &workoutRet)
		} else {
			token := c.GetHeader("Authorization")

			resp, err := shared.PositionsRequestStrWorkout(workoutRet, token)
			if err != nil {
				c.JSON(400, gin.H{
					"Error": "Issue with positions API",
					"Exact": err.Error(),
				})
				return
			}
			c.JSON(201, gin.H{
				"workout":   workoutRet,
				"positions": resp,
			})
		}
	}
}
