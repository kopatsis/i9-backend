package intro

import (
	// "fmt"

	"fulli9/intro/alteredfuncs"
	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"

	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateIntroWorkout(minutes float32, userID string, database *mongo.Database, boltDB *bbolt.DB) (shared.Workout, error) {

	user, stretches, exercises, _, typeMatrix, err := dbinput.AllInputsAsync(database, boltDB, userID)
	if err != nil {
		return shared.Workout{}, err
	}

	levelSteps := []float32{650, 700, 800, 950, 1150, 1400, 1700, 2400, 3300}

	// Uses new system
	allowedNormal, allowedCombo, allowedSplit := alteredfuncs.FilterExers(exercises, user, levelSteps)

	// Uses new system
	ratings := alteredfuncs.ExerRatings(exercises, user)

	// Uses new system
	types := alteredfuncs.SelectTypes(levelSteps, minutes)

	// Uses new system
	stretchTimes, exerTimes := alteredfuncs.CreateTimes(minutes, types)

	// Uses new system
	exerIDs := alteredfuncs.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit, exercises)
	genRatings := adjustments.GeneralTyping(exerIDs, types, exercises)

	// Uses new system
	reps, pairs := alteredfuncs.GetReps(typeMatrix, minutes, levelSteps, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics, stretchTimes, err := selections.SelectStretches(stretchTimes, stretches, 100.1, exerIDs, exercises, user.BannedStretches)
	if err != nil {
		return shared.Workout{}, err
	}

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, -1, minutes, pairs, [9]float32{}, float32(0), genRatings)
	workout.IsIntro = true

	id, err := dboutput.SaveNewWorkout(database, workout)
	if err != nil {
		return shared.Workout{}, err
	}
	workout.ID = id

	return workout, nil

}
