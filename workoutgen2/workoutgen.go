package workoutgen2

import (
	// "fmt"
	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"

	// "fulli9/workoutgen2/dbhandler"
	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"
	"fulli9/workoutgen2/stretches"

	"go.mongodb.org/mongo-driver/mongo"
	// "fulli9/workoutgen2/userinput"
)

func WorkoutGen(minutes float32, difficulty int, userID string, database *mongo.Database) (shared.AnyWorkout, error) {

	// client, database, err := dbhandler.ConnectDB()
	// if err != nil {
	// 	fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
	// 	return shared.Workout{}, err
	// }
	// defer dbhandler.DisConnectDB(client)

	// minutes, difficulty, username := userinput.GetUserInputs()

	if minutes < 10 || difficulty == 0 {
		user, err := dbinput.GetUserDB(database, userID)
		if err != nil {
			return shared.StretchWorkout{}, err
		}

		stretchWorkout, err := stretches.GetStretchWO(user, minutes, database)
		if err != nil {
			return shared.StretchWorkout{}, err
		}

		id, err := dboutput.SaveStretchWorkout(database, stretchWorkout)
		if err != nil {
			return shared.StretchWorkout{}, err
		}

		stretchWorkout.ID = id

		return stretchWorkout, nil
	}

	user, stretches, exercises, pastWOs, typeMatrix, err := dbinput.AllInputsAsync(database, userID)
	if err != nil {
		return shared.Workout{}, err
	}

	adjlevel := adjustments.CalcNewLevel(difficulty, user.Level, pastWOs)

	allowedNormal, allowedCombo, allowedSplit := adjustments.FilterExers(exercises, user, adjlevel)

	ratings := adjustments.ExerRatings(exercises, pastWOs, user)

	types := selections.SelectTypes(adjlevel, minutes)

	stretchTimes, exerTimes := creation.CreateTimes(minutes, types)

	exerIDs := selections.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit)

	reps, pairs := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics, err := selections.SelectStretches(stretchTimes, stretches, adjlevel, exerIDs, exercises, user.BannedStretches)
	if err != nil {
		return shared.Workout{}, err
	}

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, difficulty, minutes, pairs)

	id, err := dboutput.SaveNewWorkout(database, workout)
	if err != nil {
		return shared.Workout{}, err
	}
	workout.ID = id

	return workout, nil

}
