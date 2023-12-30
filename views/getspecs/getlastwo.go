package getspecs

import (
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/workoutgen2/datatypes"
	"fulli9/workoutgen2/dbhandler"
)

func GetLastWOToView(username string) datatypes.Workout {
	client, database, err := dbhandler.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
		return datatypes.Workout{}
	}
	defer dbhandler.DisConnectDB(client)
	WO := dbinput.GetPastWOsDB(database, username)
	return WO
}
