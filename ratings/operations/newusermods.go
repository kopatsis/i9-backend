package operations

import (
	"fulli9/shared"
	"math"
)

func NewUserMods(user shared.User, ratings [9]int, workout shared.Workout, exercises map[string]shared.Exercise, count int64, average int, onlyWorkout bool) (map[string]float32, map[string]float32, map[int]float32, map[int]float32) {

	exerMods := map[string]float32{}
	typeMods := map[string]float32{}

	roundEndurance := map[int]float32{}
	timeEndurance := map[int]float32{}

	if user.ExerModifications != nil {
		exerMods = user.ExerModifications
	}
	if user.TypeModifications != nil {
		typeMods = user.TypeModifications
	}
	if user.RoundEndurance != nil {
		roundEndurance = user.RoundEndurance
	}
	if user.TimeEndurance != nil {
		timeEndurance = user.TimeEndurance
	}

	if !onlyWorkout {
		for i, round := range workout.Exercises {
			for _, id := range round.ExerciseIDs {
				if val, ok := exerMods[id]; ok {
					exerMods[id] = val * (1 + (2.5*float32(workout.Difficulty-1)-float32(ratings[i]))/(25*float32(math.Log(float64(count+2)))))
					exerMods[id] = float32(math.Max(.5, math.Min(2.0, float64(exerMods[id]))))
				} else {
					exerMods[id] = (1 + (2.5*float32(workout.Difficulty-1)-float32(ratings[i]))/(25*float32(math.Log(float64(count+2)))))
				}

				if val, ok := typeMods[exercises[id].Parent]; ok {
					typeMods[exercises[id].Parent] = val * (1 + (2*float32(workout.Difficulty-1)-float32(ratings[i]))/(75*float32(math.Log(float64(count+2)))))
					typeMods[exercises[id].Parent] = float32(math.Max(.5, math.Min(2.0, float64(typeMods[exercises[id].Parent]))))

				} else {
					typeMods[exercises[id].Parent] = (1 + (2*float32(workout.Difficulty-1)-float32(ratings[i]))/(75*float32(math.Log(float64(count+2)))))
				}
			}

			if val, ok := roundEndurance[i]; ok {
				roundEndurance[i] = val * (1 + (float32(average)-float32(ratings[i]))/(75*float32(math.Log(float64(count+2)))))
				roundEndurance[i] = float32(math.Max(.5, math.Min(2.0, float64(roundEndurance[i]))))
			} else {
				roundEndurance[i] = (1 + (float32(average)-float32(ratings[i]))/(75*float32(math.Log(float64(count+2)))))
			}
		}
	}

	timeRnd := int(math.Round(float64(workout.Minutes)/5.0)) * 5
	if val, ok := timeEndurance[timeRnd]; ok {
		timeEndurance[timeRnd] = val * (1 + (2.5*float32(workout.Difficulty-1)-float32(average))/(40*float32(math.Log(float64(count+2)))))
		timeEndurance[timeRnd] = float32(math.Max(.5, math.Min(2.0, float64(timeEndurance[timeRnd]))))
	} else {
		timeEndurance[timeRnd] = (1 + (2.5*float32(workout.Difficulty-1)-float32(average))/(40*float32(math.Log(float64(count+2)))))
	}

	return exerMods, typeMods, roundEndurance, timeEndurance
}
