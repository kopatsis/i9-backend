package alteredfuncs

import (
	"fulli9/shared"
	"fulli9/workoutgen2/creation"
	"math"
)

func GetReps(matrix shared.TypeMatrix, minutes float32, levelSteps []float32, times [9]shared.ExerciseTimes, user shared.User, exerIDs [9][]string, exers map[string]shared.Exercise, types [9]string) ([9][]float32, [9][]bool) {

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
			unAdjReps := creation.UnadjustedReps(i+1, id, levelSteps[i], minutes, times[i], user, exers)
			if types[i] == "Combo" {
				adjReps := unAdjReps / float32(times[i].ComboExers)
				currentReps = append(currentReps, adjReps)
			} else {
				currentReps = append(currentReps, unAdjReps)
			}
		}

		if types[i] == "Split" {
			currentReps, retPairs = creation.SplitReps(currentReps, matrix, exers, round, switchRepTotal, retPairs, i)
		}

		for i, rep := range currentReps {
			currentReps[i] = float32(math.Max(float64(rep), 1))
		}

		retReps[i] = currentReps
	}
	return retReps, retPairs
}
