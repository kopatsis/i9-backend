package adapts

import (
	"fulli9/adapts/alteredfuncs"

	"fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"
	"fulli9/workoutgen2/dboutput"

	// "fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/mongo"
)

func Adapt(difficulty int, userID string, database *mongo.Database, workoutID string, asNew bool) (shared.Workout, error) {

	user, exercises, pastWOs, typeMatrix, workout, err := alteredfuncs.AllInputsAsync(database, userID, workoutID)
	if err != nil {
		return shared.Workout{}, nil
	}

	minutes := workout.Minutes

	adjlevel := adjustments.CalcDiffLevel(difficulty, adjustments.CalcInitLevel(user.Level, pastWOs))

	var types [9]string
	var exerTimes [9]shared.ExerciseTimes
	var exerIDs [9][]string

	for i := 0; i < 9; i++ {
		types[i] = workout.Exercises[i].Status
		exerTimes[i] = workout.Exercises[i].Times
		exerIDs[i] = workout.Exercises[i].ExerciseIDs
	}

	stretchTimes := workout.StretchTimes

	types, exerIDs, exerTimes, cardioRatings, cardioRating := adjustments.RateCardio(exerIDs, types, exerTimes, exercises, typeMatrix)
	genRatings := adjustments.GeneralTyping(exerIDs, types, exercises)

	reps, pairs := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics := workout.Statics, workout.Dynamics

	newworkout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, difficulty, minutes, pairs, cardioRatings, cardioRating, genRatings)

	if asNew {
		newworkout.Name = shared.NameAnimals(false)
		id, err := dboutput.SaveNewWorkout(database, newworkout)
		if err != nil {
			return shared.Workout{}, err
		}
		newworkout.ID = id
	} else {
		newworkout.ID = workout.ID
		newworkout.Name = workout.Name
		err := alteredfuncs.SaveUpdatedWorkout(database, newworkout)
		if err != nil {
			return shared.Workout{}, err
		}
	}

	return newworkout, nil
}
