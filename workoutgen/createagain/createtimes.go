package createagain

import (
	"math"
)

func MinutesToTimes(minutes float64) map[string]float32 {
	seconds := 60 * minutes
	ret := map[string]float32{}
	if minutes < 10 {
		return nil
	}

	var stretchSecs float32
	var dynamicRest float32
	if minutes < 20 {
		stretchSecs = float32(math.Max(1.5, seconds/8))
		dynamicRest = 10.0
	} else {
		stretchSecs = float32(math.Min(5.0, math.Max(2.5, seconds/12)))
		dynamicRest = 15.0
	}

	dynamicUsable := stretchSecs - dynamicRest

	staticSets := float32(math.Round(float64(stretchSecs / 15)))
	dynamicSets := float32(math.Round(float64(dynamicUsable / 15)))

	staticTimePer := stretchSecs / staticSets
	dynamicTimePer := dynamicUsable / dynamicSets

	exerTimeTotal := float32(seconds - (float64(stretchSecs * 2)))
	exerTimePerRound := exerTimeTotal / 9

	exerRestBetweenRounds := float32(math.Min(90, math.Max(15, 0.2*float64(exerTimePerRound))))
	exerUsableRoundTime := exerTimePerRound - exerRestBetweenRounds
	exerSetsPerRound := float32(math.Round(float64(exerTimePerRound / 30)))
	exerTimePerSet := exerUsableRoundTime / exerSetsPerRound
	exerTimePerSetReps := exerTimePerSet * (2 / 3)
	exerTimePerSetRest := exerTimePerSet * (1 / 3)

	ret["staticSecs"] = stretchSecs
	ret["staticSets"] = staticSets
	ret["staticPerSet"] = staticTimePer

	ret["dynamicSecs"] = stretchSecs
	ret["dynamicSets"] = dynamicSets
	ret["dynamicPerSet"] = dynamicTimePer
	ret["dynamicRest"] = dynamicRest

	ret["exerTotal"] = exerTimeTotal
	ret["exerPerRound"] = exerTimePerRound
	ret["exerRoundRest"] = exerRestBetweenRounds
	ret["exerPerSet"] = exerTimePerSet
	ret["exerSetsPerRound"] = exerSetsPerRound
	ret["exerTimePerSetReps"] = exerTimePerSetReps
	ret["exerTimePerSetRest"] = exerTimePerSetRest

	return ret
}
