package operations

import (
	"fulli9/workoutgen2/datatypes"
	"math"
)

func NewUserMods(user datatypes.User, ratings [9]float32, workout datatypes.Workout, exercises map[string]datatypes.Exercise, count int64) (map[string]float32, map[string]float32) {

	exerMods := user.ExerModifications
	typeMods := user.TypeModifications

	for i, round := range workout.Exercises {
		for _, id := range round.ExerciseIDs {
			if val, ok := exerMods[id]; ok {
				exerMods[id] = val * (1 + (2*float32(workout.Difficulty-1)-ratings[i])/(10*float32(math.Log(float64(count+2)))))
			} else {
				exerMods[id] = 1 * (1 + (2*float32(workout.Difficulty-1)-ratings[i])/(10*float32(math.Log(float64(count+2)))))
			}
			if val, ok := typeMods[exercises[id].Parent]; ok {
				typeMods[exercises[id].Parent] = val * (1 + (2*float32(workout.Difficulty-1)-ratings[i])/(25*float32(math.Log(float64(count+2)))))
			} else {
				typeMods[exercises[id].Parent] = 1 * (1 + (2*float32(workout.Difficulty-1)-ratings[i])/(25*float32(math.Log(float64(count+2)))))

			}
		}
	}

	return exerMods, typeMods
}
