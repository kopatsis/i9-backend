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

		c.JSON(201, &workoutRet)
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

		c.JSON(201, &workoutRet)
	}
}
