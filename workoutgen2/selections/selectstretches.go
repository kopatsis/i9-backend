package selections

import (
	"fulli9/shared"
	"fulli9/workoutgen2/stretches"
	"math/rand"
)

func SelectStretches(stretchtimes shared.StretchTimes, stretchMap map[string][]shared.Stretch, adjlevel float32, exerIDs [9][]string, exercises map[string]shared.Exercise, user shared.User) ([]string, []string, shared.StretchTimes, error) {

	bodyparts := map[int]bool{}

	for _, round := range exerIDs {
		for _, id := range round {
			for _, part := range exercises[id].BodyParts {
				bodyparts[part] = true
			}
		}
	}

	filteredStretches, err := stretches.FilterStretches(adjlevel, stretchMap, bodyparts, user)
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

		reqGroup := 0

		if stretchtimes.DynamicSets > 5 && i > stretchtimes.DynamicSets-3 {
			if i == stretchtimes.DynamicSets-1 && !stretches.ContainsReqGroup(dynamics, 1) {
				reqGroup = 1
			} else if i == stretchtimes.DynamicSets-2 && !stretches.ContainsReqGroup(dynamics, 2) {
				reqGroup = 2
			}
		}

		current := dynamicSt[rand.Intn(len(dynamicSt))]
		if reqGroup != 0 && stretches.ContainsReqGroup(dynamicSt, reqGroup) {
			count := 0

			for _, st := range dynamicSt {
				if st.ReqGroup == reqGroup {
					count++
					if rand.Intn(count) == 0 {
						current = st
					}
				}
			}
		} else {
			ct := 0
			current = stretches.SelectDynamic(dynamicSt, sum)
			for stretches.ForLoopConditions(dynamics, dynamicSt, current, ct) {
				current = stretches.SelectDynamic(dynamicSt, sum)
				ct++
			}
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
