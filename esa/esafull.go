package esa

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"fulli9/esa/fromxl"
	"fulli9/esa/mongo"
)

func RunESA() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Exercise (E), Stretch (S), or Both (B): ")
	entry, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	client, database := mongo.ConnectDB()

	entry = entry[:len(entry)-2]
	if entry != "S" && entry != "E" && entry != "B" {
		fmt.Println("Invalid Entry. Out.")
		return
	}

	if entry == "S" || entry == "B" {
		allSts := fromxl.EnterSt()

		collection := database.Collection("stretch")

		insertStretchResults := mongo.SaveStretch(collection, allSts)

		fromxl.AddStrToXL(insertStretchResults)
	}

	if entry == "E" || entry == "B" {
		allExs := fromxl.EnterEx()

		collection := database.Collection("exercise")

		insertExerciseResults := mongo.SaveExercise(collection, allExs)

		fromxl.AddExerToXL(insertExerciseResults)
	}

	client.Disconnect(context.Background())

}
