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

	statics, dynamics := []shared.Stretch{}, []shared.Stretch{}
	for i := 0; i < stretchSets; i++ {
		current := stretches["Static"][int(rand.Float64()*float64(len(stretches["Static"])))]
		for ForLoopConditions(statics, stretches["Static"], current) {
			current = stretches["Static"][int(rand.Float64()*float64(len(stretches["Static"])))]
		}
		statics = append(statics, current)

		current = stretches["Dynamic"][int(rand.Float64()*float64(len(stretches["Dynamic"])))]
		for ForLoopConditions(dynamics, stretches["Dynamic"], current) {
			current = stretches["Dynamic"][int(rand.Float64()*float64(len(stretches["Dynamic"])))]
		}
		dynamics = append(dynamics, current)
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
