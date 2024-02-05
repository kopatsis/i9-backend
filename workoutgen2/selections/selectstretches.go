package selections

import (
	"fulli9/shared"
	"fulli9/workoutgen2/stretches"
	"math/rand"
)

func SelectStretches(stretchtimes shared.StretchTimes, stretchMap map[string][]shared.Stretch, adjlevel float32, exerIDs [9][]string, exercises map[string]shared.Exercise) ([]string, []string, error) {

	bodyparts := map[int]bool{}

	for _, round := range exerIDs {
		for _, id := range round {
			for _, part := range exercises[id].BodyParts {
				bodyparts[part] = true
			}
		}
	}

	filteredStretches, err := stretches.FilterStretches(adjlevel, stretchMap, bodyparts)
	if err != nil {
		return nil, nil, err
	}

	statics, dynamics := []string{}, []string{}
	for i := 0; i < stretchtimes.StaticSets; i++ {
		statics = append(statics, filteredStretches["Static"][int(rand.Float64()*float64(len(filteredStretches["Static"])))].ID.Hex())
	}

	for i := 0; i < stretchtimes.DynamicSets; i++ {
		dynamics = append(dynamics, filteredStretches["Dynamic"][int(rand.Float64()*float64(len(filteredStretches["Dynamic"])))].ID.Hex())
	}

	return statics, dynamics, nil
}
