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
