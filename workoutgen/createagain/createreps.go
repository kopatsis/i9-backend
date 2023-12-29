package createagain

import (
	"fulli9/workoutgen/datatypes"
	"math"
	"strconv"
)

func repsForExer(useLevel float64, repMap map[int]float32) int {
	closestRep, closestLevel := -1, float32(-1.0)
	for rep, level := range repMap {
		if closestLevel == -1 {
			closestLevel = level
			closestRep = rep
		} else if math.Abs(useLevel-float64(level)) < math.Abs(useLevel-float64(closestLevel)) {
			closestLevel = level
			closestRep = rep
		}
	}
	return closestRep
}

func repsForComboEx(useLevel float64, repMap map[int]float32) int {
	closestRep, closestLevel := -1, float32(-1.0)
	for rep, level := range repMap {
		if closestLevel == -1 {
			closestLevel = level
			closestRep = rep
		} else if math.Abs(useLevel-float64(level)) < math.Abs(useLevel-float64(closestLevel)) {
			closestLevel = level
			closestRep = rep
		}
	}

	if closestRep%2 == 0 {
		closestRep = closestRep / 2
	} else {
		closestRep = closestRep/2 + 1
	}

	return closestRep
}

func repsForSplitEx(useLevel1, useLevel2 float64, repMap1, repMap2 map[int]float32, compatMod float32) int {
	reps1 := repsForExer(useLevel1, repMap1)
	reps2 := repsForExer(useLevel2, repMap2)
	average := (float32(reps1) + float32(reps2)) / 2
	return int(average * compatMod)
}

func GetReps(exercises map[string]datatypes.Exercise, user datatypes.User, exersToDo [9]map[string]string, times map[string]float32, adjlevel float64) [9]map[string]string {

	WOSpecLevel := adjlevel * float64(1+(6-times["exerSetsPerRound"])/10)

	ret := exersToDo

	for i, exerMap := range exersToDo {
		if exerMap["status"] == "COMBO" {

			exer1Level := WOSpecLevel
			if mod, ok := user.ExerSpecModifications[exerMap["id1"]]; ok {
				exer1Level *= mod
			}
			ret[i]["reps1"] = strconv.Itoa(repsForComboEx(exer1Level, exercises[exerMap["id1"]].Reps))

			exer2Level := WOSpecLevel
			if mod, ok := user.ExerSpecModifications[exerMap["id2"]]; ok {
				exer2Level *= mod
			}
			ret[i]["reps2"] = strconv.Itoa(repsForComboEx(exer2Level, exercises[exerMap["id2"]].Reps))

		} else if exerMap["status"] == "SPLIT" {

			exer1Level := WOSpecLevel
			if mod, ok := user.ExerSpecModifications[exerMap["id1"]]; ok {
				exer1Level *= mod
			}

			exer2Level := WOSpecLevel
			if mod, ok := user.ExerSpecModifications[exerMap["id2"]]; ok {
				exer2Level *= mod
			}

			ret[i]["reps"] = strconv.Itoa(repsForSplitEx(exer1Level, exer2Level, exercises[exerMap["id1"]].Reps, exercises[exerMap["id2"]].Reps, exercises[exerMap["id1"]].Compatibles[exerMap["id2"]][0]))
		} else {
			exerLevel := WOSpecLevel
			if mod, ok := user.ExerSpecModifications[exerMap["id"]]; ok {
				exerLevel *= mod
			}
			ret[i]["reps"] = strconv.Itoa(repsForExer(exerLevel, exercises[exerMap["id"]].Reps))
		}
	}

	return ret
}
