package workoutgen

import (
	"fmt"
	"fulli9/workoutgen/adjustments"
	"fulli9/workoutgen/createagain"
	"fulli9/workoutgen/dbhandler"
	"fulli9/workoutgen/inputs"
	"fulli9/workoutgen/stretches"
)

func FullGenerate() {

	client, database, err := dbhandler.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
		return
	}
	defer dbhandler.DisConnectDB(client)

	minutes, difficulty, username := inputs.GetUserInputs()
	user, pastWOs := inputs.GetDBInputs(database, username)
	if difficulty == 0 || minutes < 10 {
		fmt.Println(stretches.GetStretchWO(user.Level, minutes, database))
	}

	adjLevel := adjustments.CalcNewLevel(difficulty, user.Level, pastWOs)

	allExercises := inputs.GetAllExers(database)
	filteredExercises := adjustments.FilterExers(allExercises, user, adjLevel)
	allRatings := adjustments.ExerRatings(filteredExercises, pastWOs, adjLevel)
	timesMap := createagain.MinutesToTimes(minutes)

	exersToDo := createagain.CreateExerList(filteredExercises, allRatings)
	exersToDo = createagain.GetReps(filteredExercises, user, exersToDo, timesMap, adjLevel)
	stretchMap := stretches.GetStretches(database, filteredExercises, timesMap, exersToDo, adjLevel)

	finalWO := createagain.ConstructWO(exersToDo, timesMap, stretchMap, user.ID.Hex(), user.Level)

	fmt.Println(finalWO)

}
