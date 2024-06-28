package stretches

import (
	"fulli9/shared"
	"fulli9/workoutgen2/dbinput"
	"math"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func StretchTimeSlice(strList []shared.Stretch, usableTime float32) []float32 {
	ret := []float32{}

	var sum float32
	sum = 0

	for _, str := range strList {
		if str.InPairs {
			sum += 1.667
		} else {
			sum++
		}
	}

	unpairedTime := usableTime / sum
	pairedTime := unpairedTime * 1.667

	for _, str := range strList {
		if str.InPairs {
			ret = append(ret, pairedTime)
		} else {
			ret = append(ret, unpairedTime)
		}
	}

	return ret

}

func StretchToString(strList []shared.Stretch) []string {
	ret := []string{}
	for _, str := range strList {
		ret = append(ret, str.ID.Hex())
	}
	return ret
}

func GetStretchWO(user shared.User, minutes float32, database *mongo.Database) (shared.StretchWorkout, error) {
	stretches, err := dbinput.GetStretchesDB(database)
	if err != nil {
		return shared.StretchWorkout{}, err
	}

	stretches, err = FilterStretches(user.Level*1.3, stretches, nil, user.BannedStretches)
	if err != nil {
		return shared.StretchWorkout{}, err
	}

	secsPerSet, circles := 15.0, 1
	if minutes > 2 && minutes < 8 {
		secsPerSet = 20
	} else if minutes < 16 {
		secsPerSet = 20
		circles = 2
	} else if minutes < 24 {
		secsPerSet = 20
		circles = 3
	} else if minutes < 36 {
		secsPerSet = 25
		circles = 3
	} else if minutes < 48 {
		secsPerSet = 30
		circles = 3
	} else if minutes < 72 {
		secsPerSet = 30
		circles = 4
	} else {
		secsPerSet = 30
		circles = 5
	}

	stretchSecs := (60 * minutes) / 2
	stretchSecsCircled := stretchSecs / float32(circles)
	stretchSets := int(math.Round(float64(stretchSecsCircled) / secsPerSet))

	dynamicSt := stretches["Dynamic"]
	staticSt := stretches["Static"]

	sum := float32(0)
	for _, st := range dynamicSt {
		sum += st.Weight
	}

	statics, dynamics := []shared.Stretch{}, []shared.Stretch{}
	for i := 0; i < stretchSets; i++ {

		reqGroup := 0

		if stretchSets > 3 && i > stretchSets-3 {
			if i == stretchSets-1 && !ContainsReqGroup(dynamics, 1) {
				reqGroup = 1
			} else if i == stretchSets-2 && !ContainsReqGroup(dynamics, 2) {
				reqGroup = 2
			}
		}

		current := dynamicSt[rand.Intn(len(dynamicSt))]
		if reqGroup != 0 && ContainsReqGroup(dynamicSt, reqGroup) {
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
			current := SelectDynamic(dynamicSt, sum)
			for ForLoopConditions(dynamics, dynamicSt, current) {
				current = SelectDynamic(dynamicSt, sum)

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

	realstatics, realdynamics := []shared.Stretch{}, []shared.Stretch{}
	for i := 0; i < circles; i++ {
		realstatics = append(realstatics, statics...)
		realdynamics = append(realdynamics, dynamics...)
	}

	ret := shared.StretchWorkout{
		Name:   "",
		UserID: user.ID.Hex(),
		Date:   primitive.NewDateTimeFromTime(time.Now()),
		Status: "Not Started",
		StretchTimes: shared.StretchTimes{
			DynamicPerSet: StretchTimeSlice(realdynamics, stretchSecs),
			StaticPerSet:  StretchTimeSlice(realstatics, stretchSecs),
			DynamicSets:   stretchSets,
			StaticSets:    stretchSets,
			DynamicRest:   0.0,
			FullRound:     stretchSecs,
		},
		LevelAtStart: user.Level,
		Dynamics:     StretchToString(realdynamics),
		Statics:      StretchToString(realstatics),
	}

	return ret, nil

}

func SelectDynamic(dynamics []shared.Stretch, sum float32) shared.Stretch {
	randSelect := rand.Float32() * sum
	for _, st := range dynamics {
		randSelect -= st.Weight
		if randSelect <= 0.1 {
			return st
		}
	}
	return dynamics[int(rand.Float64()*float64(len(dynamics)))]

}

func ContainsReqGroup(stlist []shared.Stretch, group int) bool {
	for _, st := range stlist {
		if st.ReqGroup == 1 {
			return true
		}
	}
	return false
}

func ForLoopConditions(existing, filtered []shared.Stretch, current shared.Stretch) bool {
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
