package ratings

import (
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostRating(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var rateHandler shared.RateRoute
		if err := c.ShouldBindJSON(&rateHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		ratings := [9]float32{-1, -1, -1, -1, -1, -1, -1, -1, -1}
		for i, score := range rateHandler.Ratings {
			if i > 9 {
				break
			}
			ratings[i] = score
		}

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		if err := RateWorkout(rateHandler.UserID, ratings, id, database); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with rating route",
				"Exact": err.Error(),
			})
			return
		}

		c.Status(204)
	}
}

func PostIntroRating(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var rateHandler shared.RateIntroRoute
		if err := c.ShouldBindJSON(&rateHandler); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		if err := RateIntroWorkout(rateHandler.UserID, rateHandler.Rounds, database); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with intro rating route",
				"Exact": err.Error(),
			})
			return
		}

		c.Status(204)
	}
}
