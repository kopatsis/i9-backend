package selections

import (
	"fulli9/workoutgen2/datatypes"
	"fulli9/workoutgen2/stretches"
	"math/rand"
)

func SelectStretches(stretchtimes datatypes.StretchTimes, stretchMap map[string][]datatypes.Stretch, adjlevel float32, exerIDs [9][]string, exercises map[string]datatypes.Exercise) ([]string, []string) {

	bodyparts := map[int]bool{}

	for _, round := range exerIDs {
		for _, id := range round {
			for _, part := range exercises[id].BodyParts {
				bodyparts[part] = true
			}
		}
	}

	filteredStretches := stretches.FilterStretches(adjlevel, stretchMap, bodyparts)

	statics, dynamics := []string{}, []string{}
	for i := 0; i < stretchtimes.StaticSets; i++ {
		statics = append(statics, filteredStretches["Static"][int(rand.Float64()*float64(len(filteredStretches["Static"])))].ID.Hex())
	}

	for i := 0; i < stretchtimes.DynamicSets; i++ {
		dynamics = append(dynamics, filteredStretches["Dynamic"][int(rand.Float64()*float64(len(filteredStretches["Dynamic"])))].ID.Hex())
	}

	return statics, dynamics
}
