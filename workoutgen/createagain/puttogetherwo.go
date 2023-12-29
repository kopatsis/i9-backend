package createagain

import (
	"fulli9/workoutgen/datatypes"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConstructWO(exercises [9]map[string]string, times map[string]float32, stretches map[string][]string, id string, level float64) datatypes.Workout {
	var finalWO datatypes.Workout

	finalWO.UserID = id
	finalWO.Created = primitive.NewDateTimeFromTime(time.Now())
	finalWO.Status = "Unstarted"

	finalWO.Statics = stretches["statics"]
	finalWO.Dynamics = stretches["dynamics"]
	finalWO.Times = times

	finalWO.LevelAtStart = level
	finalWO.Exercises = []map[string]any{}
	for _, exer := range exercises {
		finalWO.Exercises = append(finalWO.Exercises, exer)
	}

	return finalWO
}
