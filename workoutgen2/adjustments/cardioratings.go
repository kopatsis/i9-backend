package adjustments

import (
	"fulli9/shared"
	"sort"
)

type indexedValue struct {
	value float32
	index int
}

func RateCardio(exerIDs [9][]string, types [9]string, exerTimes [9]shared.ExerciseTimes, exercises map[string]shared.Exercise, typeMatrix shared.TypeMatrix) ([9]string, [9][]string, [9]shared.ExerciseTimes, [9]float32, float32) {

	parentMatIndex := map[string]int{
		"Pushups":           0,
		"Squats":            1,
		"Burpees":           2,
		"Jumps":             3,
		"Lunges":            4,
		"Mountain Climbers": 5,
		"Abs":               6,
		"Bridges":           7,
		"Kicks":             8,
		"Planks":            9,
		"Supermans":         10,
	}

	cardioRatings, cardioRating := [9]float32{}, float32(0)

	for i, idlist := range exerIDs {
		if types[i] == "Regular" {
			cardioRatings[i] = exercises[idlist[0]].CardioRating
		} else if types[i] == "Combo" {
			var sum float32
			var count float32
			for _, id := range idlist {
				sum += exercises[id].CardioRating
				count += 1
			}
			cardioRatings[i] = (sum / count) * 1.4
		} else {
			exer1, exer2 := exercises[idlist[0]], exercises[idlist[1]]
			sum := exer1.CardioRating + exer2.CardioRating

			typeMultiplier := 1 / typeMatrix.Matrix[parentMatIndex[exer1.Parent]][parentMatIndex[exer2.Parent]]
			cardioRatings[i] = (sum / 2) * (((1.6 * 3) + typeMultiplier*2) / 5)
		}

		cardioRating += cardioRatings[i]
	}

	posArray := sortedIndicesByPattern(cardioRatings)

	newCardioRatings := specialSortFloat32(cardioRatings, posArray)
	newTypes := specialSortString(types, posArray)
	newExerIDs := specialSortSlice(exerIDs, posArray)
	newExerTimes := specialSortExerTimes(exerTimes, posArray)

	cardioRating /= 9

	return newTypes, newExerIDs, newExerTimes, newCardioRatings, cardioRating
}

func sortedIndicesByPattern(arr [9]float32) [9]int {

	pattern := [9]int{1, 3, 5, 7, 8, 9, 6, 4, 2}

	indexedArr := make([]indexedValue, len(arr))
	for i, v := range arr {
		indexedArr[i] = indexedValue{v, i}
	}

	sort.Slice(indexedArr, func(i, j int) bool {
		return indexedArr[i].value < indexedArr[j].value
	})

	var result [9]int
	for i, p := range pattern {
		result[i] = indexedArr[p-1].index
	}

	return result
}

func specialSortFloat32(originalArray [9]float32, indexOrder [9]int) [9]float32 {
	var sortedArray [9]float32
	for i, index := range indexOrder {
		sortedArray[i] = originalArray[index]
	}
	return sortedArray
}

func specialSortString(originalArray [9]string, indexOrder [9]int) [9]string {
	var sortedArray [9]string
	for i, index := range indexOrder {
		sortedArray[i] = originalArray[index]
	}
	return sortedArray
}

func specialSortSlice(originalArray [9][]string, indexOrder [9]int) [9][]string {
	var sortedArray [9][]string
	for i, index := range indexOrder {
		sortedArray[i] = originalArray[index]
	}
	return sortedArray
}

func specialSortExerTimes(originalArray [9]shared.ExerciseTimes, indexOrder [9]int) [9]shared.ExerciseTimes {
	restAdj := [9]float32{0.85, 0.9, 0.95, 1, 1.1, 1.1, 1.15, 1.1, 0.85}
	var sortedArray [9]shared.ExerciseTimes
	for i, index := range indexOrder {
		sortedArray[i] = originalArray[index]
		oldrest := sortedArray[i].RestPerRound
		newrest := oldrest * restAdj[i]
		newFT := (sortedArray[i].FullRound - oldrest) + newrest
		sortedArray[i].FullRound = newFT
		sortedArray[i].RestPerRound = newrest
	}
	return sortedArray
}
