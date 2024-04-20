package intro

import (
	// "fmt"

	"fulli9/intro/alteredfuncs"
	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"

	// "fulli9/workoutgen2/dbhandler"
	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"

	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateIntroWorkout(minutes float32, userID string, database *mongo.Database) (shared.Workout, error) {

	user, stretches, exercises, _, typeMatrix, err := dbinput.AllInputsAsync(database, userID)
	if err != nil {
		return shared.Workout{}, err
	}
	pastWOs := []shared.Workout{}

	levelSteps := []float32{50, 125, 200, 350, 500, 800, 1100, 1700, 2300}

	// Uses new system
	allowedNormal, allowedCombo, allowedSplit := alteredfuncs.FilterExers(exercises, user, levelSteps)

	ratings := adjustments.ExerRatings(5, exercises, pastWOs, user)

	// Uses new system
	types := alteredfuncs.SelectTypes(levelSteps, minutes)

	stretchTimes, exerTimes := creation.CreateTimes(minutes, types)

	// Uses new system
	exerIDs := alteredfuncs.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit, exercises)
	genRatings := adjustments.GeneralTyping(exerIDs, types, exercises)

	// Uses new system
	reps, pairs := alteredfuncs.GetReps(typeMatrix, minutes, levelSteps, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics, stretchTimes, err := selections.SelectStretches(stretchTimes, stretches, levelSteps[0], exerIDs, exercises, user.BannedStretches)
	if err != nil {
		return shared.Workout{}, err
	}

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, -1, minutes, pairs, [9]float32{}, float32(0), genRatings)

	id, err := dboutput.SaveNewWorkout(database, workout)
	if err != nil {
		return shared.Workout{}, err
	}
	workout.ID = id

	return workout, nil

}
