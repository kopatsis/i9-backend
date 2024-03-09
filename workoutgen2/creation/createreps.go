package creation

import (
	"fulli9/shared"
	"math"
)

func unadjustedReps(round int, id string, adjlevel, minutes float32, times shared.ExerciseTimes, user shared.User, exers map[string]shared.Exercise) float32 {

	exercise := exers[id]

	speclevel := adjlevel

	if val, ok := user.TimeEndurance[int(math.Round(float64(minutes)/5.0))*5]; ok {
		speclevel *= val
	}
	if val, ok := user.RoundEndurance[round]; ok {
		speclevel *= val
	}
	if val, ok := user.ExerModifications[id]; ok {
		speclevel *= val
	}
	if val, ok := user.TypeModifications[exercise.Parent]; ok {
		speclevel *= val
	}

	// speclevel *= 1 + (float32(40-minutes) / 160)

	varA := exercise.RepVars[0]
	varB := exercise.RepVars[1]
	varC := exercise.RepVars[2]

	speclevel = float32(math.Max(float64(speclevel), float64(exercise.MinLevel)))

	initReps := float32(exercise.MinReps) + varA*float32(math.Pow(float64((speclevel)/varB), 1/float64(varC)))

	initReps *= (times.ExercisePerSet / 20)

	// Wouldn't ever be < 5, but don't want a crash
	if minutes >= 5 {
		initReps *= float32(2/math.Log2(float64(minutes+1))) - .275
	}

	return initReps
}

func splitReps(currentReps []float32, matrix shared.TypeMatrix, exers map[string]shared.Exercise, round []string, switchRepTotal float32, retPairs [9][]bool, i int) ([]float32, [9][]bool) {

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
		"Planks":            9,
		"MISC":              10,
	}

	exer1 := exers[round[0]]
	exer2 := exers[round[1]]

	if exer1.InPairs && exer2.InPairs {
		switchRepTotal = float32(math.Max(2.0, float64(switchRepTotal)*.5))
		retPairs[i][0] = true
		retPairs[i][1] = true
	} else if exer1.InPairs {
		switchRepTotal = float32(math.Max(2.0, float64(switchRepTotal)*.75))
		retPairs[i][0] = true
	} else if exer2.InPairs {
		switchRepTotal = float32(math.Max(2.0, float64(switchRepTotal)*.75))
		retPairs[i][1] = true
	}

	repadjust := matrix.Matrix[parentMatIndex[exer1.Parent]][parentMatIndex[exer2.Parent]]

	adjReps := (switchRepTotal * repadjust) / 2

	currentReps = append(currentReps, adjReps)

	return currentReps, retPairs
}

func GetReps(matrix shared.TypeMatrix, minutes, adjlevel float32, times [9]shared.ExerciseTimes, user shared.User, exerIDs [9][]string, exers map[string]shared.Exercise, types [9]string) ([9][]float32, [9][]bool) {

	retReps := [9][]float32{}
	retPairs := [9][]bool{}
	for i, ids := range exerIDs {
		add := []bool{}
		for range ids {
			add = append(add, false)
		}
		retPairs[i] = add
	}

	for i, round := range exerIDs {
		currentReps := []float32{}
		switchRepTotal := float32(0)
		for _, id := range round {
			unAdjReps := unadjustedReps(i+1, id, adjlevel, minutes, times[i], user, exers)
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
			currentReps, retPairs = splitReps(currentReps, matrix, exers, round, switchRepTotal, retPairs, i)
		}
		retReps[i] = currentReps
	}
	return retReps, retPairs
}
