package adjustments

import (
	"fulli9/workoutgen/datatypes"
	"math"
	"time"
)

func ExerRatings(exercises map[string]datatypes.Exercise, pastWOs []datatypes.Workout, adjlevel float64) map[string]float32 {
	ret := map[string]float32{}
	for _, exercise := range exercises {
		if exercise.Name == "SPLIT" {
			ret[exercise.ID.Hex()] = float32(exercise.StartQuality) + float32(adjlevel/4)
		} else if exercise.Name == "COMBO" {
			ret[exercise.ID.Hex()] = float32(exercise.StartQuality) + float32(adjlevel/10)
		} else {
			ret[exercise.ID.Hex()] = float32(exercise.StartQuality)
		}
	}
	for _, workout := range pastWOs {
		adjustment := int(time.Now().Sub(workout.Created.Time()).Hours())
		for _, exercise := range workout.Exercises {
			ret[exercise["ID"].(string)] = float32(math.Max(0.0, float64(ret[exercise["ID"].(string)]-float32(11-adjustment))))
		}
	}
	return ret
}
