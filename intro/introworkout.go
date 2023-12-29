package intro

import (
	"fmt"
	"fulli9/intro/alteredfuncs"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"
	"fulli9/workoutgen2/datatypes"
	"fulli9/workoutgen2/dbhandler"
	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"
)

func GenerateIntroWorkout(username string) (datatypes.Workout, error) {

	client, database, err := dbhandler.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
		return datatypes.Workout{}, err
	}
	defer dbhandler.DisConnectDB(client)

	// user := dbinput.GetUserDB(database, username)
	// stretches := dbinput.GetStetchesDB(database)
	// exercises := dbinput.GetExersDB(database)
	// pastWOs := []datatypes.Workout{}
	// typeMatrix := dbinput.GetMatrix(database)

	user, stretches, exercises, _, typeMatrix := dbinput.AllInputsAsync(database, username)
	pastWOs := []datatypes.Workout{}

	minutes := float32(45)

	levelSteps := []float32{50, 125, 225, 375, 575, 825, 1125, 1475, 1975}

	// Uses new system
	allowedNormal, allowedCombo, allowedSplit := alteredfuncs.FilterExers(exercises, user, levelSteps)

	ratings := adjustments.ExerRatings(exercises, pastWOs, user)

	// Uses new system
	types := alteredfuncs.SelectTypes(levelSteps, minutes)

	stretchTimes, exerTimes := creation.CreateTimes(minutes, types)

	// Uses new system
	exerIDs := alteredfuncs.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit)

	// Uses new system
	reps := alteredfuncs.GetReps(typeMatrix, minutes, levelSteps, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics := selections.SelectStretches(stretchTimes, stretches, levelSteps[0], exerIDs, exercises)

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, -1)

	dboutput.SaveNewWorkout(database, workout)

	return workout, nil

}
