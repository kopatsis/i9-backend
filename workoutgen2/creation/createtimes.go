package creation

import (
	"fulli9/shared"
	"math"
	"math/rand"
)

func CreateTimes(minutes float32, types [9]string) (shared.StretchTimes, [9]shared.ExerciseTimes) {
	possibleComboTimes := [3]float32{30.0, 45.0, 60.0}
	possibleSplitTimes := [3]float32{30.0, 45.0, 30.0}

	seconds := 60 * minutes

	var retStr shared.StretchTimes
	retExer := [9]shared.ExerciseTimes{}

	if minutes < 20 {
		retStr.FullRound = float32(math.Max(1.5, float64(seconds/8)))
		retStr.DynamicRest = 10.0
	} else {
		retStr.FullRound = float32(math.Min(5.0, math.Max(2.5, float64(seconds/12))))
		retStr.DynamicRest = 15.0
	}

	dynamicUsable := retStr.FullRound - retStr.DynamicRest

	staticSets := float32(math.Round(float64(retStr.FullRound / 15)))
	dynamicSets := float32(math.Round(float64(dynamicUsable / 15)))

	retStr.StaticPerSet = retStr.FullRound / staticSets
	retStr.DynamicPerSet = dynamicUsable / dynamicSets

	retStr.StaticSets = int(staticSets)
	retStr.DynamicSets = int(dynamicSets)

	exerTimeTotal := float32(seconds - (retStr.FullRound * 2))
	exerTimePerRound := exerTimeTotal / 9

	exerRestBetweenRounds := float32(math.Min(90, math.Max(15, 0.2*float64(exerTimePerRound))))
	exerUsableRoundTime := exerTimePerRound - exerRestBetweenRounds

	for i, round := range types {
		var currentTimes shared.ExerciseTimes

		var roundSets float32

		if round == "Combo" && exerUsableRoundTime >= 2 {
			comboNum := int(rand.Float64() * 3)
			timeC := possibleComboTimes[comboNum]
			roundSets = float32(math.Round(float64(exerUsableRoundTime / timeC)))
			currentTimes.ComboExers = comboNum + 2
		} else if round == "Split" && exerUsableRoundTime >= 2 {
			splitNum := int(rand.Float64() * 3)
			timeS := possibleSplitTimes[splitNum]
			roundSets = float32(math.Round(float64(exerUsableRoundTime / timeS)))
		} else {
			roundSets = float32(math.Round(float64(exerUsableRoundTime / 30)))
			if round == "Combo" {
				currentTimes.ComboExers = 2
			} else {
				currentTimes.ComboExers = 0
			}
		}
		exerTimePerSet := exerUsableRoundTime / roundSets
		currentTimes.ExercisePerSet = (exerTimePerSet * 2) / 3
		currentTimes.RestPerSet = (exerTimePerSet * 1) / 3
		currentTimes.RestPerRound = exerRestBetweenRounds
		currentTimes.FullRound = exerTimePerRound
		currentTimes.Sets = int(roundSets)

		retExer[i] = currentTimes
	}

	return retStr, retExer
}
