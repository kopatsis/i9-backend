package alteredfuncs

import (
	"fulli9/shared"
	"math"
	"math/rand"
)

func CreateTimes(minutes float32, types [9]string) (shared.StretchTimes, [9]shared.ExerciseTimes) {
	possibleComboTimes := [5]float32{30.0, 30.0, 30.0, 45.0, 45.0}
	possibleSplitTimes := [5]float32{30.0, 45.0, 30.0, 45.0, 45.0}

	seconds := 60 * minutes

	var retStr shared.StretchTimes
	retExer := [9]shared.ExerciseTimes{}

	if minutes < 20 {
		retStr.FullRound = float32(math.Max(1.5*60, float64(seconds/8)))
		retStr.DynamicRest = 10.0
	} else {
		retStr.FullRound = float32(math.Min(5.0*60, math.Max(2.5*60, float64(seconds/12))))
		retStr.DynamicRest = 15.0
	}

	dynamicUsable := retStr.FullRound - retStr.DynamicRest

	secsPerSetInit := 20.0
	if retStr.FullRound < 75 {
		secsPerSetInit = 12
	} else if retStr.FullRound < 120 {
		secsPerSetInit = 16
	}

	dynamicSets := float32(math.Round(float64(dynamicUsable) / secsPerSetInit))

	retStr.StaticSets = int(dynamicSets)
	retStr.DynamicSets = int(dynamicSets)

	exerTimeTotal := float32(seconds - (retStr.FullRound * 2))
	exerTimePerRound := exerTimeTotal / 9

	exerRestBetweenRounds := float32(math.Min(90, math.Max(15, 0.2*float64(exerTimePerRound))))
	if minutes < 20 {
		exerRestBetweenRounds = float32(math.Min(90, math.Max(15, 0.125*float64(exerTimePerRound))))
	} else if minutes < 30 {
		exerRestBetweenRounds = float32(math.Min(90, math.Max(15, 0.15*float64(exerTimePerRound))))
	}
	exerUsableRoundTime := exerTimePerRound - exerRestBetweenRounds

	restAdj := [9]float32{0.85, 0.9, 0.95, 1, 1.1, 1.1, 1.15, 1.1, 0.85}

	for i, round := range types {
		var currentTimes shared.ExerciseTimes

		var roundSets float32

		if round == "Combo" {
			comboNum := int(rand.Float64() * 5)
			timeC := possibleComboTimes[comboNum]
			roundSets = float32(math.Round(float64(exerUsableRoundTime / timeC)))
			currentTimes.ComboExers = 2
		} else if round == "Split" {
			splitNum := int(rand.Float64() * 5)
			timeS := possibleSplitTimes[splitNum]
			roundSets = float32(math.Round(float64(exerUsableRoundTime / timeS)))
		} else {
			roundSets = float32(math.Round(float64(exerUsableRoundTime / 30)))
		}

		exerTimePerSet := exerUsableRoundTime / roundSets
		currentTimes.ExercisePerSet = (exerTimePerSet * 2) / 3
		currentTimes.RestPerSet = (exerTimePerSet * 1) / 3
		currentTimes.RestPerRound = exerRestBetweenRounds * restAdj[i]
		currentTimes.FullRound = (exerTimePerRound - exerRestBetweenRounds) + currentTimes.RestPerRound
		currentTimes.Sets = int(roundSets)

		retExer[i] = currentTimes
	}

	return retStr, retExer
}
