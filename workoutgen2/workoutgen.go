package workoutgen2

import (
	"fmt"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"
	"fulli9/workoutgen2/datatypes"
	"fulli9/workoutgen2/dbhandler"
	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"
	"fulli9/workoutgen2/stretches"
	// "fulli9/workoutgen2/userinput"
)

func WorkoutGen(minutes float32, difficulty int, username string) (datatypes.AnyWorkout, error) {

	client, database, err := dbhandler.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
		return datatypes.Workout{}, err
	}
	defer dbhandler.DisConnectDB(client)

	// minutes, difficulty, username := userinput.GetUserInputs()

	if minutes < 10 || difficulty == 0 {
		user := dbinput.GetUserDB(database, username)
		stretchWorkout := stretches.GetStretchWO(user, minutes, database)
		dboutput.SaveStretchWorkout(database, stretchWorkout)
		return stretchWorkout, nil
	}

	user, stretches, exercises, pastWOs, typeMatrix := dbinput.AllInputsAsync(database, username)

	adjlevel := adjustments.CalcNewLevel(difficulty, user.Level, pastWOs)

	allowedNormal, allowedCombo, allowedSplit := adjustments.FilterExers(exercises, user, adjlevel)

	ratings := adjustments.ExerRatings(exercises, pastWOs, user)

	types := selections.SelectTypes(adjlevel, minutes)

	stretchTimes, exerTimes := creation.CreateTimes(minutes, types)

	exerIDs := selections.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit)

	reps := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics := selections.SelectStretches(stretchTimes, stretches, adjlevel, exerIDs, exercises)

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, difficulty)

	dboutput.SaveNewWorkout(database, workout)

	return workout, nil

}
