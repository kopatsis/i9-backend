package ratings

import (
	"fulli9/shared"
	"math"

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

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		ratings := [9]float32{-1, -1, -1, -1, -1, -1, -1, -1, -1}
		for i, score := range rateHandler.Ratings {
			if i > 9 {
				break
			}
			ratings[i] = float32(math.Min(math.Max(float64(score), 10), 0))
		}

		favorites := [9]float32{-1, -1, -1, -1, -1, -1, -1, -1, -1}
		for i, score := range rateHandler.Favoritism {
			if i > 9 {
				break
			}
			favorites[i] = float32(math.Min(math.Max(float64(score), 10), 0))
		}

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL parameter",
			})
			return
		}

		if err := RateWorkout(userID, ratings, favorites, rateHandler.FullRating, rateHandler.FullFave, rateHandler.OnlyWorkout, id, database); err != nil {
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

		userID, err := shared.GetIDFromReq(database, c)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with userID",
				"Exact": err.Error(),
			})
			return
		}

		if err := RateIntroWorkout(userID, rateHandler.Rounds, database); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with intro rating route",
				"Exact": err.Error(),
			})
			return
		}

		c.Status(204)
	}
}
