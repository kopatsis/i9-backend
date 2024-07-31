package ratings

import (
	"context"
	"fmt"
	"fulli9/shared"
	"math"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostQuiz(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var results shared.QuizRoute
		if err := c.ShouldBindJSON(&results); err != nil {
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

		if err := PushQuiz(results, database, userID); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with saving user result",
				"Exact": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"Message": "Success",
		})

	}
}

func PushQuiz(results shared.QuizRoute, database *mongo.Database, userID string) error {
	level := float32(0)
	pushups := "Wall"
	plyo := 0

	switch results.Stamina {
	case 0:
		level = 300
	case 1:
		level = 100
	case 2:
		level = 300
	case 3:
		level = 600
	case 4:
		level = 1000
	}

	switch results.Endurance {
	case 0:
		level *= 1
	case 1:
		level *= 0.667
	case 2:
		level *= 0.925
	case 3:
		level *= 1.25
	case 4:
		level *= 1.75
	}

	switch results.LowerStrength {
	case 1:
		level = float32(math.Max(math.Min(float64(level-100), float64(level*2)), float64(level*0.667)))
	case 3:
		level = float32(math.Max(math.Min(float64(level+200), float64(level*2)), float64(level*0.667)))
	case 4:
		level = float32(math.Max(math.Min(float64(level+400), float64(level*2)), float64(level*0.667)))
	}

	switch results.UpperStrength {
	case 0:
		pushups = "Knee"
	case 1:
		pushups = "Wall"
	case 2:
		pushups = "Knee"
	case 3:
		pushups = "Regular"
	case 4:
		pushups = "Regular"
		plyo += 1
	}

	switch results.Plyo {
	case 0:
		plyo += 1
	case 2:
		plyo += 1
	case 3:
		plyo += 2
	case 4:
		plyo += 3
	}

	if err := UpdateUser(database, userID, level, plyo, pushups); err != nil {
		return err
	}
	return nil
}

func UpdateUser(database *mongo.Database, userID string, level float32, plyo int, pushup string) error {
	collection := database.Collection("user")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"level":       level,
			"plyoToler":   plyo,
			"pushsetting": pushup,
			"assessed":    true,
		},
		"$inc": bson.M{
			"displevel": 8,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}
