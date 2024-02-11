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

	// client, database, err := dbhandler.ConnectDB()
	// if err != nil {
	// 	fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
	// 	return shared.Workout{}, err
	// }
	// defer dbhandler.DisConnectDB(client)

	// user := dbinput.GetUserDB(database, username)
	// stretches := dbinput.GetStetchesDB(database)
	// exercises := dbinput.GetExersDB(database)
	// pastWOs := []shared.Workout{}
	// typeMatrix := dbinput.GetMatrix(database)

	user, stretches, exercises, _, typeMatrix, err := dbinput.AllInputsAsync(database, userID)
	if err != nil {
		return shared.Workout{}, err
	}
	pastWOs := []shared.Workout{}

	// minutes := float32(45)

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

	statics, dynamics, err := selections.SelectStretches(stretchTimes, stretches, levelSteps[0], exerIDs, exercises, user.BannedStretches)
	if err != nil {
		return shared.Workout{}, err
	}

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, -1, minutes)

	dboutput.SaveNewWorkout(database, workout)

	return workout, nil

}
