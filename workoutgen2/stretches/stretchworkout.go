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
			sum += 2
		} else {
			sum++
		}
	}

	unpairedTime := usableTime / sum
	pairedTime := unpairedTime * 2

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

	stretches, err = FilterStretches(user.Level*1.1, stretches, nil, user.BannedStretches)
	if err != nil {
		return shared.StretchWorkout{}, err
	}

	stretchSecs := (60 * minutes) / 2

	secsPerSetInit := 20.0
	if minutes < 2 {
		secsPerSetInit = 15
	} else if minutes > 8 {
		secsPerSetInit = 30
	} else if minutes > 24 {
		secsPerSetInit = 40
	}

	stretchSets := int(math.Round(float64(stretchSecs) / secsPerSetInit))

	statics, dynamics := []shared.Stretch{}, []shared.Stretch{}
	for i := 0; i < stretchSets; i++ {
		current := stretches["Static"][int(rand.Float64()*float64(len(stretches["Static"])))]
		for len(statics) > 0 && statics[len(statics)-1].ID.Hex() == current.ID.Hex() && len(stretches["Static"]) > 1 {
			current = stretches["Static"][int(rand.Float64()*float64(len(stretches["Static"])))]
		}
		statics = append(statics, current)

		current = stretches["Dynamic"][int(rand.Float64()*float64(len(stretches["Dynamic"])))]
		for len(dynamics) > 0 && dynamics[len(dynamics)-1].ID.Hex() == current.ID.Hex() && len(stretches["Dynamic"]) > 1 {
			current = stretches["Dynamic"][int(rand.Float64()*float64(len(stretches["Dynamic"])))]
		}
		dynamics = append(dynamics, current)

	}

	ret := shared.StretchWorkout{
		Name:   "",
		UserID: user.ID.Hex(),
		Date:   primitive.NewDateTimeFromTime(time.Now()),
		Status: "Not Started",
		StretchTimes: shared.StretchTimes{
			DynamicPerSet: StretchTimeSlice(dynamics, stretchSecs),
			StaticPerSet:  StretchTimeSlice(statics, stretchSecs),
			DynamicSets:   stretchSets,
			StaticSets:    stretchSets,
			DynamicRest:   0.0,
			FullRound:     stretchSecs,
		},
		LevelAtStart: user.Level,
		Dynamics:     StretchToString(dynamics),
		Statics:      StretchToString(statics),
	}

	return ret, nil

}
