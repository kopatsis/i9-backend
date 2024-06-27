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

	dynamicSt := filteredStretches["Dynamic"]
	staticSt := filteredStretches["Static"]

	sum := float32(0)
	for _, st := range dynamicSt {
		sum += st.Weight
	}

	statics, dynamics := []shared.Stretch{}, []shared.Stretch{}

	for i := 0; i < stretchtimes.DynamicSets; i++ {
		current := stretches.SelectDynamic(dynamicSt, sum)
		for stretches.ForLoopConditions(dynamics, dynamicSt, current) {
			current = stretches.SelectDynamic(dynamicSt, sum)

		}
		dynamics = append(dynamics, current)

		staticID := current.DynamicPairs[int(rand.Float64()*float64(len(current.DynamicPairs)))]
		currentStatic := shared.Stretch{}
		for _, str := range staticSt {
			if str.ID.Hex() == staticID {
				currentStatic = str
			}
		}
		if currentStatic.Name == "" {
			currentStatic = staticSt[int(rand.Float64()*float64(len(staticSt)))]
		}
		statics = append(statics, currentStatic)

	}

	retStretchTimes := stretchtimes
	retStretchTimes.DynamicPerSet = stretches.StretchTimeSlice(dynamics, retStretchTimes.FullRound-retStretchTimes.DynamicRest)
	retStretchTimes.StaticPerSet = stretches.StretchTimeSlice(statics, retStretchTimes.FullRound)

	return stretches.StretchToString(statics), stretches.StretchToString(dynamics), retStretchTimes, nil
}
