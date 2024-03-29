package intro

import (
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostIntroWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.IntroWorkoutRoute
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
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		workout, err := GenerateIntroWorkout(woHandler.Time, userID, database)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with workout generator",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(201, &workout)
	}
}
