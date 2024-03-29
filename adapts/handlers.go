package adapts

import (
	"fulli9/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostAdaptedWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.AdaptWorkoutRoute
		woHandler.AsNew = true

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

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		workout, err := Adapt(woHandler.Difficulty, userID, database, id, woHandler.AsNew)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with adapted workout generator",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(201, &workout)
	}
}

func PostExternalAdaptedWorkout(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var woHandler shared.AdaptWorkoutRoute
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

		id, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL paramete",
			})
			return
		}

		workout, err := Adapt(woHandler.Difficulty, userID, database, id, true)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with adapted workout generator",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(201, &workout)
	}
}
