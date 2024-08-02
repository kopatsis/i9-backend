package workoutgen2

import (
	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"

	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"
	"fulli9/workoutgen2/stretches"

	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/mongo"
)

func WorkoutGen(minutes float32, difficulty int, loweronly bool, userID string, database *mongo.Database, boltDB *bbolt.DB) (shared.User, shared.AnyWorkout, error) {

	if minutes < 8 || difficulty == 0 {
		return stretchWOReturn(minutes, userID, database, boltDB)
	}

	user, stretches, exercises, pastWOs, typeMatrix, err := dbinput.AllInputsAsync(database, boltDB, userID)
	if err != nil {
		return shared.User{}, shared.Workout{}, err
	}

	adjlevel := adjustments.CalcInitLevel(user.Level, pastWOs)

	allowedNormal, allowedCombo, allowedSplit := adjustments.FilterExers(difficulty, exercises, user, adjlevel, loweronly)

	adjlevel = adjustments.CalcDiffLevel(difficulty, adjlevel)

	ratings := adjustments.ExerRatings(difficulty, loweronly, exercises, pastWOs, user)
	adjustments.AdjustBurpeeRatings(user, exercises, ratings)

	types := selections.SelectTypes(adjlevel, minutes, difficulty)

	stretchTimes, exerTimes := creation.CreateTimes(minutes, types)

	exerIDs := selections.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit, exercises)

	types, exerIDs, exerTimes, cardioRatings, cardioRating := adjustments.RateCardio(exerIDs, types, exerTimes, exercises, typeMatrix)
	genRatings := adjustments.GeneralTyping(exerIDs, types, exercises)

	reps, pairs := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics, stretchTimes, err := selections.SelectStretches(stretchTimes, stretches, adjlevel, exerIDs, exercises, user)
	if err != nil {
		return shared.User{}, shared.Workout{}, err
	}

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, difficulty, minutes, pairs, cardioRatings, cardioRating, genRatings)

	id, err := dboutput.SaveNewWorkout(database, workout)
	if err != nil {
		return shared.User{}, shared.Workout{}, err
	}
	workout.ID = id
	workout.LowerOnly = loweronly

	// No longer doing this as last minutes is just default minutes now
	// err = dboutput.UpdateUserLast(minutes, difficulty, userID, database)
	// if err != nil {
	// 	return shared.Workout{}, err
	// }

	return user, workout, nil
}

func stretchWOReturn(minutes float32, userID string, database *mongo.Database, boltDB *bbolt.DB) (shared.User, shared.StretchWorkout, error) {
	user, err := dbinput.GetUserDB(database, userID)
	if err != nil {
		return shared.User{}, shared.StretchWorkout{}, err
	}

	stretchWorkout, err := stretches.GetStretchWO(user, minutes, database, boltDB)
	if err != nil {
		return shared.User{}, shared.StretchWorkout{}, err
	}
	stretchWorkout.Minutes = minutes

	stretchWorkout.Name = shared.NameAnimals(true)
	id, err := dboutput.SaveStretchWorkout(database, stretchWorkout)
	if err != nil {
		return shared.User{}, shared.StretchWorkout{}, err
	}

	stretchWorkout.ID = id

	return user, stretchWorkout, nil
}
