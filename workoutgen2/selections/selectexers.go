package selections

import (
	"fulli9/shared"
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
	return exers[int(rand.Float32()*float32(len(exers)))]
}

func getSplitIDs(sum float32, exers []string, ratings map[string]float32, exercises map[string]shared.Exercise, ret [9][]string) []string {
	id1 := selectID(sum, exers, ratings)

	for !AllowedByGroup(ret, id1, exercises) {
		id1 = selectID(sum, exers, ratings)
	}

	id2 := selectID(sum, exers, ratings)
	for id1 == id2 || exercises[id1].Parent == exercises[id2].Parent || !AllowedByGroup(ret, id2, exercises) {
		id2 = selectID(sum, exers, ratings)
	}

	return []string{id1, id2}
}

func getComboIDs(count int, sum float32, exers []string, ratings map[string]float32, exercises map[string]shared.Exercise, existing [9][]string) []string {
	ret := []string{}
	for i := 0; i < count; i++ {
		currentID := selectID(sum, exers, ratings)
		for slices.Contains(ret, currentID) || (i > 0 && exercises[ret[i-1]].Parent == exercises[currentID].Parent) || !AllowedByGroup(existing, currentID, exercises) {
			currentID = selectID(sum, exers, ratings)
		}
		ret = append(ret, currentID)
	}
	return ret
}

func checkPrevRoundsByType(rounds [9][]string, types [9]string, current []string, i int) bool {
	for j := 0; j < i; j++ {
		if types[i] == types[j] && areSlicesSame(rounds[j], current) {
			return false
		}
	}
	return true
}

func areSlicesSame(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func SelectExercises(types [9]string, times [9]shared.ExerciseTimes, ratings map[string]float32, allowedNormal, allowedCombo, allowedSplit []string, exercises map[string]shared.Exercise) [9][]string {
	normalSum := sum(allowedNormal, ratings)
	comboSum := sum(allowedCombo, ratings)
	splitSum := sum(allowedSplit, ratings)

	ret := [9][]string{}

	for i, round := range types {
		if round == "Combo" {
			current := getComboIDs(times[i].ComboExers, comboSum, allowedCombo, ratings, exercises, ret)
			for !checkPrevRoundsByType(ret, types, current, i) {
				current = getComboIDs(times[i].ComboExers, comboSum, allowedCombo, ratings, exercises, ret)
			}
			ret[i] = current
		} else if round == "Split" {
			current := getSplitIDs(splitSum, allowedSplit, ratings, exercises, ret)
			for !checkPrevRoundsByType(ret, types, current, i) {
				current = getSplitIDs(splitSum, allowedSplit, ratings, exercises, ret)
			}
			ret[i] = current
		} else {
			current := []string{selectID(normalSum, allowedNormal, ratings)}
			for !checkPrevRoundsByType(ret, types, current, i) || !AllowedByGroup(ret, current[0], exercises) {
				current = []string{selectID(normalSum, allowedNormal, ratings)}
			}
			ret[i] = current
		}
	}

	return ret
}

func AllowedByGroup(ret [9][]string, id string, exercises map[string]shared.Exercise) bool {
	if exercises[id].SinglesGroup == 0 {
		return true
	}
	for _, rnd := range ret {
		for _, rid := range rnd {
			if exercises[id].SinglesGroup == exercises[rid].SinglesGroup {
				return false
			}
		}
	}
	return true
}
