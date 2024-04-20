package workoutgen2

import (
	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"

	"fulli9/workoutgen2/dbinput"
	"fulli9/workoutgen2/dboutput"
	"fulli9/workoutgen2/selections"
	"fulli9/workoutgen2/stretches"

	"go.mongodb.org/mongo-driver/mongo"
)

func WorkoutGen(minutes float32, difficulty int, userID string, database *mongo.Database) (shared.AnyWorkout, error) {

	if minutes < 8 || difficulty == 0 {
		return stretchWOReturn(minutes, userID, database)
	}

	user, stretches, exercises, pastWOs, typeMatrix, err := dbinput.AllInputsAsync(database, userID)
	if err != nil {
		return shared.Workout{}, err
	}

	adjlevel := adjustments.CalcNewLevel(difficulty, user.Level, pastWOs)

	allowedNormal, allowedCombo, allowedSplit := adjustments.FilterExers(exercises, user, adjlevel)

	ratings := adjustments.ExerRatings(exercises, pastWOs, user)
	adjustments.AdjustBurpeeRatings(user, exercises, ratings)

	types := selections.SelectTypes(adjlevel, minutes, difficulty)

	stretchTimes, exerTimes := creation.CreateTimes(minutes, types)

	exerIDs := selections.SelectExercises(types, exerTimes, ratings, allowedNormal, allowedCombo, allowedSplit)

	types, exerIDs, exerTimes, cardioRatings, cardioRating := adjustments.RateCardio(exerIDs, types, exerTimes, exercises, typeMatrix)
	genRatings := adjustments.GeneralTyping(exerIDs, types, exercises)

	reps, pairs := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics, stretchTimes, err := selections.SelectStretches(stretchTimes, stretches, adjlevel, exerIDs, exercises, user.BannedStretches)
	if err != nil {
		return shared.Workout{}, err
	}

	workout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, difficulty, minutes, pairs, cardioRatings, cardioRating, genRatings)

	id, err := dboutput.SaveNewWorkout(database, workout)
	if err != nil {
		return shared.Workout{}, err
	}
	workout.ID = id

	err = dboutput.UpdateUserLast(minutes, difficulty, userID, database)
	if err != nil {
		return shared.Workout{}, err
	}

	return workout, nil
}

func stretchWOReturn(minutes float32, userID string, database *mongo.Database) (shared.StretchWorkout, error) {
	user, err := dbinput.GetUserDB(database, userID)
	if err != nil {
		return shared.StretchWorkout{}, err
	}

	stretchWorkout, err := stretches.GetStretchWO(user, minutes, database)
	if err != nil {
		return shared.StretchWorkout{}, err
	}
	stretchWorkout.Minutes = minutes

	stretchWorkout.Name = shared.NameAnimals(true)
	id, err := dboutput.SaveStretchWorkout(database, stretchWorkout)
	if err != nil {
		return shared.StretchWorkout{}, err
	}

	stretchWorkout.ID = id

	return stretchWorkout, nil
}
