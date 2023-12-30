package views

import (
	"encoding/json"
	"fmt"
	"fulli9/views/getspecs"
)

func ViewLastWorkout(username string) {
	PrettyPrint(getspecs.GetLastWOToView(username))
}

func PrettyPrint(obj interface{}) {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	fmt.Println(string(bytes))
}
