package alteredfuncs

import (
	"fulli9/shared"
	"math/rand"
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
	id1 := "6607a0c0147215fc5694bb5b"

	for id, ex := range exercises {
		if ex.Name == "Jump Lunges" {
			id1 = id
		}
	}

	id2 := selectID(sum, exers, ratings)

	if rand.Float32() >= 0.5 {
		return []string{id1, id2}
	} else {
		return []string{id2, id1}
	}

}

func getComboIDs(count int, sum float32, exers []string, ratings map[string]float32, exercises map[string]shared.Exercise) []string {
	id1 := "6607a0bf147215fc5694bb4a"

	for id, ex := range exercises {
		if ex.Name == "High Knees" {
			id1 = id
		}
	}

	id2 := selectID(sum, exers, ratings)

	if rand.Float32() >= 0.5 {
		return []string{id1, id2}
	} else {
		return []string{id2, id1}
	}
}

func SelectExercises(types [9]string, times [9]shared.ExerciseTimes, ratings map[string]float32, allowedNormal, allowedCombo, allowedSplit [9][]string, exercises map[string]shared.Exercise) [9][]string {

	ret := [9][]string{}

	for i, round := range types {

		if round == "Combo" {
			sum := sum(allowedCombo[i], ratings)
			current := getComboIDs(times[i].ComboExers, sum, allowedCombo[i], ratings, exercises)
			ret[i] = current
		} else if round == "Split" {
			sum := sum(allowedSplit[i], ratings)
			current := getSplitIDs(sum, allowedSplit[i], ratings, exercises)
			ret[i] = current
		} else {
			sum := sum(allowedNormal[i], ratings)
			current := []string{selectID(sum, allowedNormal[i], ratings)}
			ret[i] = current
		}
	}

	return ret
}
