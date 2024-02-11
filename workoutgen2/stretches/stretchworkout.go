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

	stretchSets := int(math.Round(float64(stretchSecs) / 20))

	stretchSecsPerSet := stretchSecs / float32(stretchSets)

	statics, dynamics := []string{}, []string{}
	for i := 0; i < stretchSets; i++ {
		statics = append(statics, stretches["Static"][int(rand.Float64()*float64(len(stretches["Static"])))].ID.Hex())
		dynamics = append(dynamics, stretches["Dynamic"][int(rand.Float64()*float64(len(stretches["Dynamic"])))].ID.Hex())
	}

	ret := shared.StretchWorkout{
		Name:   "",
		UserID: user.ID.Hex(),
		Date:   primitive.NewDateTimeFromTime(time.Now()),
		Status: "Not Started",
		StretchTimes: shared.StretchTimes{
			DynamicPerSet: stretchSecsPerSet,
			StaticPerSet:  stretchSecsPerSet,
			DynamicSets:   stretchSets,
			StaticSets:    stretchSets,
			DynamicRest:   0.0,
			FullRound:     stretchSecs,
		},
		LevelAtStart: user.Level,
		Dynamics:     dynamics,
		Statics:      statics,
	}

	return ret, nil

}
