package selections

import (
	"fulli9/shared"
	"fulli9/workoutgen2/stretches"
	"math/rand"
)

func SelectStretches(stretchtimes shared.StretchTimes, stretchMap map[string][]shared.Stretch, adjlevel float32, exerIDs [9][]string, exercises map[string]shared.Exercise, bannedStretches []string) ([]string, []string, shared.StretchTimes, error) {

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
		return nil, nil, shared.StretchTimes{}, err
	}

	statics, dynamics := []shared.Stretch{}, []shared.Stretch{}
	for i := 0; i < stretchtimes.StaticSets; i++ {
		current := filteredStretches["Static"][int(rand.Float64()*float64(len(filteredStretches["Static"])))]
		for forLoopConditions(statics, filteredStretches["Static"], current) {
			current = filteredStretches["Static"][int(rand.Float64()*float64(len(filteredStretches["Static"])))]
		}
		statics = append(statics, current)
	}

	for i := 0; i < stretchtimes.DynamicSets; i++ {
		current := filteredStretches["Dynamic"][int(rand.Float64()*float64(len(filteredStretches["Dynamic"])))]
		for forLoopConditions(dynamics, filteredStretches["Dynamic"], current) {
			current = filteredStretches["Dynamic"][int(rand.Float64()*float64(len(filteredStretches["Dynamic"])))]
		}
		dynamics = append(dynamics, current)
	}

	retStretchTimes := stretchtimes
	retStretchTimes.DynamicPerSet = stretches.StretchTimeSlice(dynamics, retStretchTimes.FullRound-retStretchTimes.DynamicRest)
	retStretchTimes.StaticPerSet = stretches.StretchTimeSlice(statics, retStretchTimes.FullRound)

	return stretches.StretchToString(statics), stretches.StretchToString(dynamics), retStretchTimes, nil
}

func forLoopConditions(existing, filtered []shared.Stretch, current shared.Stretch) bool {
	if len(existing) == 0 || len(filtered) < 2 {
		return false
	}

	if existing[len(existing)-1].ID.Hex() == current.ID.Hex() {
		return true
	}

	if len(existing) > 1 && len(filtered) > 2 && existing[len(existing)-2].ID.Hex() == current.ID.Hex() {
		return true
	}

	if len(existing) > 2 && len(filtered) > 3 && existing[len(existing)-3].ID.Hex() == current.ID.Hex() {
		return true
	}

	if len(existing) <= len(filtered) {
		count := 0
		for _, stretch := range existing {
			if stretch.ID.Hex() == current.ID.Hex() {
				count++
				if count >= 2 {
					return true
				}
			}
		}
	}

	return false
}
