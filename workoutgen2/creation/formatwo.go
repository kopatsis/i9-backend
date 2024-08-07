package creation

import (
	"fulli9/shared"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FormatWorkout(statics, dynamics []string, reps [9][]float32, exerIDs [9][]string, stretchTimes shared.StretchTimes, exerTimes [9]shared.ExerciseTimes, types [9]string, user shared.User, difficulty int, minutes float32, pairs [9][]bool, cardioRatings [9]float32, cardioRating float32, genRatings [3]float32) shared.Workout {
	ret := shared.Workout{
		Name:            shared.NameAnimals(false),
		UserID:          user.ID.Hex(),
		Username:        user.Username,
		Created:         primitive.NewDateTimeFromTime(time.Now()),
		Minutes:         minutes,
		Status:          "Unstarted",
		StretchTimes:    stretchTimes,
		LevelAtStart:    user.Level,
		Difficulty:      difficulty,
		Dynamics:        dynamics,
		Statics:         statics,
		CardioRatings:   cardioRatings,
		CardioRating:    cardioRating,
		GeneralTypeVals: genRatings,
		AvgRating:       float32(-1),
		AvgFaves:        float32(-1),
	}

	roundSlice := [9]shared.WorkoutRound{}
	for i, idlist := range exerIDs {
		round := shared.WorkoutRound{
			ExerciseIDs: idlist,
			Reps:        reps[i],
			Pairs:       pairs[i],
			Status:      types[i],
			Times:       exerTimes[i],
			AvgRating:   float32(-1),
			AvgFaves:    float32(-1),
		}
		roundSlice[i] = round
	}
	ret.Exercises = roundSlice

	return ret
}
