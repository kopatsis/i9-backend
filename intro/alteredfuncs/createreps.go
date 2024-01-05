package alteredfuncs

import (
	"fulli9/workoutgen2/datatypes"
	"math"
)

func unadjustedReps(round int, id string, adjlevel, minutes float32, times datatypes.ExerciseTimes, user datatypes.User, exers map[string]datatypes.Exercise) float32 {

	exercise := exers[id]

	speclevel := adjlevel

	speclevel *= user.TimeEndurance[int(math.Round(float64(minutes)/5.0))*5]
	speclevel *= user.RoundEndurance[round]
	speclevel *= user.ExerModifications[id]
	speclevel *= user.TypeModifications[exercise.Parent]
	speclevel *= 1 + (float32(6-times.Sets) / 12)

	varA := exercise.RepVars[0]
	varB := exercise.RepVars[1]
	varC := exercise.RepVars[2]

	speclevel = float32(math.Max(float64(speclevel), float64(exercise.MinLevel)))

	initReps := float32(exercise.MinReps) + varA*float32(math.Pow(float64((speclevel-exercise.MinLevel)/varB), float64(varC)))

	initReps *= (times.ExercisePerSet / 30)
	return initReps
}

func GetReps(matrix datatypes.TypeMatrix, minutes float32, levelSteps []float32, times [9]datatypes.ExerciseTimes, user datatypes.User, exerIDs [9][]string, exers map[string]datatypes.Exercise, types [9]string) [9][]float32 {

	parentMatIndex := map[string]int{
		"Pushups":           0,
		"Squats":            1,
		"Burpees":           2,
		"Jumps":             3,
		"Lunges":            4,
		"Mountain Climbers": 5,
		"Crunches":          6,
		"Bridges":           7,
		"Kicks":             8,
		"MISC":              9,
	}

	ret := [9][]float32{}
	for i, round := range exerIDs {
		currentReps := []float32{}
		switchRepTotal := float32(0)
		for _, id := range round {
			unAdjReps := unadjustedReps(i+1, id, levelSteps[i], minutes, times[i], user, exers)
			if types[i] == "Combo" {
				adjReps := unAdjReps / float32(times[i].ComboExers)
				currentReps = append(currentReps, adjReps)
			} else if types[i] == "Split" {
				switchRepTotal += unAdjReps
			} else {
				currentReps = append(currentReps, unAdjReps)
			}
		}

		if types[i] == "Split" {
			exer1 := exers[round[0]]
			exer2 := exers[round[1]]

			repadjust := matrix.Matrix[parentMatIndex[exer1.Parent]][parentMatIndex[exer2.Parent]]

			adjReps := (switchRepTotal * repadjust) / 2
			currentReps = append(currentReps, adjReps)
		}
		ret[i] = currentReps
	}
	return ret
}
