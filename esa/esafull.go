package esa

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"fulli9/esa/fromxl"
	"fulli9/esa/mongodb"
)

func RunESA() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Exercise (E), Stretch (S), Both (B), or Matrix(M): ")
	entry, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	entry = entry[:len(entry)-2]
	if entry != "S" && entry != "E" && entry != "B" && entry != "M" {
		fmt.Println("Invalid Entry. Out.")
		return
	}

	if entry == "M" {
		AddMatrixToDBFull()
	}

	client, database := mongodb.ConnectDB()

	if entry == "S" || entry == "B" {
		allSts := fromxl.EnterSt()

		collection := database.Collection("stretch")

		insertStretchResults := mongodb.SaveStretch(collection, allSts)

		fromxl.AddStrToXL(insertStretchResults)
	}

	if entry == "E" || entry == "B" {
		allExs := fromxl.EnterEx()

		collection := database.Collection("exercise")

		insertExerciseResults := mongodb.SaveExercise(collection, allExs)

		fromxl.AddExerToXL(insertExerciseResults)
	}

	client.Disconnect(context.Background())

}
