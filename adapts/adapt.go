package adapts

import (
	"errors"
	"fulli9/adapts/alteredfuncs"

	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"
	"fulli9/workoutgen2/dboutput"

	"go.mongodb.org/mongo-driver/mongo"
)

func Adapt(difficulty int, userID string, database *mongo.Database, workoutID string, asNew, external bool) (shared.Workout, error) {

	user, exercises, pastWOs, typeMatrix, workout, err := alteredfuncs.AllInputsAsync(database, userID, workoutID)
	if err != nil {
		return shared.Workout{}, nil
	}

	if !external && workout.UserID != userID {
		return shared.Workout{}, errors.New("workout does not match the user that's requesting it")
	}

	minutes := workout.Minutes

	var newdiff int
	if difficulty == -1 {
		newdiff = workout.Difficulty
	} else {
		newdiff = difficulty
	}

	adjlevel := adjustments.CalcDiffLevel(newdiff, adjustments.CalcInitLevel(user.Level, pastWOs))

	var types [9]string
	var exerTimes [9]shared.ExerciseTimes
	var exerIDs [9][]string

	for i := 0; i < 9; i++ {
		types[i] = workout.Exercises[i].Status
		exerTimes[i] = workout.Exercises[i].Times
		exerIDs[i] = workout.Exercises[i].ExerciseIDs
	}

	types, exerIDs, exerTimes, cardioRatings, cardioRating := adjustments.RateCardio(exerIDs, types, exerTimes, exercises, typeMatrix)
	genRatings := adjustments.GeneralTyping(exerIDs, types, exercises)

	reps, pairs := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	if asNew {
		newworkout := creation.FormatWorkout(workout.Statics, workout.Dynamics, reps, exerIDs, workout.StretchTimes, exerTimes, types, user, difficulty, minutes, pairs, cardioRatings, cardioRating, genRatings)
		newworkout.Name = shared.NameAnimals(false)
		id, err := dboutput.SaveNewWorkout(database, newworkout)
		if err != nil {
			return shared.Workout{}, err
		}
		newworkout.ID = id
		return newworkout, nil
	} else {
		workout = alteredfuncs.ModifyExisting(workout, reps, pairs, genRatings, cardioRatings, cardioRating)
		err := alteredfuncs.SaveUpdatedWorkout(database, workout)
		if err != nil {
			return shared.Workout{}, err
		}
		return workout, nil
	}

}
