package getspecs

import (
	"fmt"
	"fulli9/ratings/dbinput"
	"fulli9/shared"
	"fulli9/workoutgen2/dbhandler"
)

func GetLastWOToView(username string) shared.Workout {
	client, database, err := dbhandler.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database %s, restart.\n", err.Error())
		return shared.Workout{}
	}
	defer dbhandler.DisConnectDB(client)
	WO, err := dbinput.GetPastWOsDB(database, username)
	if err != nil {
		fmt.Println(err)
	}
	return WO
}
