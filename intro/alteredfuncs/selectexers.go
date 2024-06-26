package alteredfuncs

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

func getSplitIDs(sum float32, exers []string, ratings map[string]float32, exercises map[string]shared.Exercise) []string {
	id1 := selectID(sum, exers, ratings)

	id2 := selectID(sum, exers, ratings)
	for id1 == id2 || exercises[id1].Parent == exercises[id2].Parent {
		id2 = selectID(sum, exers, ratings)
	}

	return []string{id1, id2}
}

func getComboIDs(count int, sum float32, exers []string, ratings map[string]float32, exercises map[string]shared.Exercise) []string {
	ret := []string{}
	for i := 0; i < count; i++ {
		currentID := selectID(sum, exers, ratings)
		for slices.Contains(ret, currentID) || (i > 0 && exercises[ret[i-1]].Parent == exercises[currentID].Parent) {
			currentID = selectID(sum, exers, ratings)
		}
		ret = append(ret, currentID)
	}
	return ret
}

func checkPrevRoundsByType(rounds [9][]string, types [9]string, current []string, i int, exercises map[string]shared.Exercise) bool {

	for _, id := range current {
		if i < 5 && i != 2 && i != 4 && exercises[id].Parent == "Pushups" {
			return false
		}
	}

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

func SelectExercises(types [9]string, times [9]shared.ExerciseTimes, ratings map[string]float32, allowedNormal, allowedCombo, allowedSplit [9][]string, exercises map[string]shared.Exercise) [9][]string {

	ret := [9][]string{}

	for i, round := range types {

		if round == "Combo" {

			if i == 4 {
				current := []string{}
				for k, v := range exercises {
					if v.Name == "Pushups" {
						current = append(current, k)
					}
				}

				sum := sum(allowedCombo[i], ratings)
				currentAdd := getComboIDs(times[i].ComboExers-1, sum, allowedCombo[i], ratings, exercises)
				for !checkPrevRoundsByType(ret, types, currentAdd, i, exercises) {
					currentAdd = getComboIDs(times[i].ComboExers-1, sum, allowedCombo[i], ratings, exercises)
				}

				current = append(current, currentAdd...)
				ret[i] = current
			} else if i == 2 {
				current := []string{}
				for k, v := range exercises {
					if v.Name == "Knee Pushups" {
						current = append(current, k)
					}
				}

				sum := sum(allowedCombo[i], ratings)
				currentAdd := getComboIDs(times[i].ComboExers-1, sum, allowedCombo[i], ratings, exercises)
				for !checkPrevRoundsByType(ret, types, currentAdd, i, exercises) {
					currentAdd = getComboIDs(times[i].ComboExers-1, sum, allowedCombo[i], ratings, exercises)
				}

				current = append(current, currentAdd...)
				ret[i] = current
			} else {
				sum := sum(allowedCombo[i], ratings)
				current := getComboIDs(times[i].ComboExers, sum, allowedCombo[i], ratings, exercises)
				for !checkPrevRoundsByType(ret, types, current, i, exercises) {
					current = getComboIDs(times[i].ComboExers, sum, allowedCombo[i], ratings, exercises)
				}
				ret[i] = current
			}
		} else if round == "Split" {
			sum := sum(allowedSplit[i], ratings)
			current := getSplitIDs(sum, allowedSplit[i], ratings, exercises)
			for !checkPrevRoundsByType(ret, types, current, i, exercises) {
				current = getSplitIDs(sum, allowedSplit[i], ratings, exercises)
			}
			ret[i] = current
		} else {
			sum := sum(allowedNormal[i], ratings)
			current := []string{selectID(sum, allowedNormal[i], ratings)}
			for !checkPrevRoundsByType(ret, types, current, i, exercises) {
				current = []string{selectID(sum, allowedNormal[i], ratings)}
			}
			ret[i] = current
		}
	}

	return ret
}
