package alteredfuncs

import (
	"fulli9/workoutgen2/datatypes"
	"math/rand"
	"slices"
)

func sum(ex []string, ratings map[string]float32) float32 {
	sum := float32(0.0)
	for _, id := range ex {
		sum += ratings[id]
	}
	return sum
}

func selectID(sum float32, exers []string, ratings map[string]float32) string {
	randWeight := rand.Float32() * sum
	for _, id := range exers {
		randWeight -= ratings[id]
		if randWeight <= 0.1 {
			return id
		}
	}
	return ""
}

func getSplitIDs(sum float32, exers []string, ratings map[string]float32) []string {
	id1 := selectID(sum, exers, ratings)

	id2 := selectID(sum, exers, ratings)
	for id1 == id2 {
		id2 = selectID(sum, exers, ratings)
	}

	return []string{id1, id2}
}

func getComboIDs(count int, sum float32, exers []string, ratings map[string]float32) []string {
	ret := []string{}
	for i := 0; i < count; i++ {
		currentID := selectID(sum, exers, ratings)
		for slices.Contains(ret, currentID) {
			currentID = selectID(sum, exers, ratings)
		}
		ret = append(ret, currentID)
	}
	return ret
}

func SelectExercises(types [9]string, times [9]datatypes.ExerciseTimes, ratings map[string]float32, allowedNormal, allowedCombo, allowedSplit [9][]string) [9][]string {

	ret := [9][]string{}

	for i, round := range types {

		if round == "Combo" {
			sum := sum(allowedCombo[i], ratings)
			ret[i] = getComboIDs(times[i].ComboExers, sum, allowedCombo[i], ratings)
		} else if round == "Split" {
			sum := sum(allowedSplit[i], ratings)
			ret[i] = getSplitIDs(sum, allowedSplit[i], ratings)
		} else {
			sum := sum(allowedNormal[i], ratings)
			ret[i] = []string{selectID(sum, allowedNormal[i], ratings)}
		}
	}

	return ret
}
