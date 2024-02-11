package selections

import (
	"fulli9/shared"
	"fulli9/workoutgen2/stretches"
	"math/rand"
)

func SelectStretches(stretchtimes shared.StretchTimes, stretchMap map[string][]shared.Stretch, adjlevel float32, exerIDs [9][]string, exercises map[string]shared.Exercise, bannedStretches []string) ([]string, []string, error) {

	bodyparts := map[int]bool{}

	for _, round := range exerIDs {
		for _, id := range round {
			for _, part := range exercises[id].BodyParts {
				bodyparts[part] = true
			}
		}
	}

	filteredStretches, err := stretches.FilterStretches(adjlevel, stretchMap, bodyparts, bannedStretches)
	if err != nil {
		return nil, nil, err
	}

	statics, dynamics := []string{}, []string{}
	for i := 0; i < stretchtimes.StaticSets; i++ {
		current := filteredStretches["Static"][int(rand.Float64()*float64(len(filteredStretches["Static"])))].ID.Hex()
		for len(statics) > 0 && statics[len(statics)-1] == current && len(filteredStretches["Static"]) > 1 {
			current = filteredStretches["Static"][int(rand.Float64()*float64(len(filteredStretches["Static"])))].ID.Hex()
		}
		statics = append(statics, current)
	}

	for i := 0; i < stretchtimes.DynamicSets; i++ {
		current := filteredStretches["Dynamic"][int(rand.Float64()*float64(len(filteredStretches["Dynamic"])))].ID.Hex()
		for len(dynamics) > 0 && dynamics[len(dynamics)-1] == current && len(filteredStretches["Dynamic"]) > 1 {
			current = filteredStretches["Dynamic"][int(rand.Float64()*float64(len(filteredStretches["Dynamic"])))].ID.Hex()
		}
		dynamics = append(dynamics, current)
	}

	return statics, dynamics, nil
}
