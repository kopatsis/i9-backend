package operations

import (
	"fulli9/shared"
	"math"
)

func NewExerciseFactorialVars(ratings [9]float32, workout shared.Workout, exercises map[string]shared.Exercise, count int64) map[string][3]float32 {
	ret := map[string][3]float32{}

	for i, round := range workout.Exercises {
		for _, id := range round.ExerciseIDs {

			currentVars := exercises[id].RepVars
			if vars, ok := ret[id]; ok {
				currentVars = vars
			}

			currentVars[0] = currentVars[0] * (1 + (2*float32(workout.Difficulty-1)-ratings[i])/(10*float32(math.Log(2*float64(count+2)))))
			currentVars[2] = currentVars[2] * (1 + (ratings[i]-2*float32(workout.Difficulty-1))/(25*float32(math.Log(2*float64(count+2)))))

			ret[id] = currentVars

		}
	}

	return ret
}
