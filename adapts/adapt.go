package adapts

import (
	"fulli9/adapts/alteredfuncs"

	// "fulli9/shared"
	"fulli9/workoutgen2/adjustments"
	"fulli9/workoutgen2/creation"
	"fulli9/workoutgen2/datatypes"

	"go.mongodb.org/mongo-driver/mongo"
)

func Adapt(difficulty int, username string, database *mongo.Database, workoutID string) (datatypes.Workout, error) {

	user, exercises, pastWOs, typeMatrix, workout := alteredfuncs.AllInputsAsync(database, username, workoutID)

	minutes := workout.Minutes

	adjlevel := adjustments.CalcNewLevel(difficulty, user.Level, pastWOs)

	var types [9]string
	var exerTimes [9]datatypes.ExerciseTimes
	var exerIDs [9][]string

	for i := 0; i < 9; i++ {
		types[i] = workout.Exercises[i].Status
		exerTimes[i] = workout.Exercises[i].Times
		exerIDs[i] = workout.Exercises[i].ExerciseIDs
	}

	stretchTimes := workout.StretchTimes

	reps := creation.GetReps(typeMatrix, minutes, adjlevel, exerTimes, user, exerIDs, exercises, types)

	statics, dynamics := workout.Statics, workout.Dynamics

	newworkout := creation.FormatWorkout(statics, dynamics, reps, exerIDs, stretchTimes, exerTimes, types, user, difficulty, minutes)

	return newworkout, nil
}
