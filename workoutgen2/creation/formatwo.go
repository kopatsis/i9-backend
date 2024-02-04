package creation

import (
	"fulli9/workoutgen2/datatypes"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FormatWorkout(statics, dynamics []string, reps [9][]float32, exerIDs [9][]string, stretchTimes datatypes.StretchTimes, exerTimes [9]datatypes.ExerciseTimes, types [9]string, user datatypes.User, difficulty int, minutes float32) datatypes.Workout {
	ret := datatypes.Workout{
		Name:         "",
		UserID:       user.ID.Hex(),
		Username:     user.Username,
		Date:         primitive.NewDateTimeFromTime(time.Now()),
		Minutes:      minutes,
		Status:       "Unstarted",
		StretchTimes: stretchTimes,
		LevelAtStart: user.Level,
		Difficulty:   difficulty,
		Dynamics:     dynamics,
		Statics:      statics,
	}

	roundSlice := [9]datatypes.WorkoutRound{}
	for i, idlist := range exerIDs {
		round := datatypes.WorkoutRound{
			ExerciseIDs: idlist,
			Reps:        reps[i],
			Status:      types[i],
			Times:       exerTimes[i],
			Rating:      float32(-1),
		}
		roundSlice[i] = round
	}
	ret.Exercises = roundSlice

	return ret
}
