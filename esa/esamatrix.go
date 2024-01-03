package esa

import (
	"context"
	"fmt"
	"fulli9/esa/mongodb"
	"fulli9/workoutgen2/datatypes"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddMatrixToDBFull() {

	client, database := mongodb.ConnectDB()
	defer client.Disconnect(context.Background())

	f, err := excelize.OpenFile("esa/i9mat.xlsx")
	if err != nil {
		fmt.Printf("Error opening matrix file: %s\n", err.Error())
		return
	}

	var matrix [10][10]float32
	for i := 2; i < 12; i++ {
		var row [10]float32
		col := string('B' + i - 2)
		for j := 2; j < 12; j++ {
			value, err := f.GetCellValue("Matrix", col+strconv.Itoa(j))
			if err != nil {
				fmt.Printf("Error with matrix on cell %s: %s\n", col+strconv.Itoa(j), err.Error())
				return
			}
			valueFloat, err := strconv.ParseFloat(value, 32)
			if err != nil {
				fmt.Printf("Error with matrix on cell %s: %s\n", col+strconv.Itoa(j), err.Error())
				return
			}
			row[j-2] = float32(valueFloat)
		}
		matrix[i-2] = row
	}

	collection := database.Collection("typematrix")

	filter := bson.D{}

	var oldmatrix datatypes.TypeMatrix
	finderr := collection.FindOne(context.Background(), filter).Decode(&oldmatrix)
	if finderr != nil && finderr != mongo.ErrNoDocuments {
		fmt.Printf("Error getting matrix from DB: %s\n", err.Error())
		return
	}

	if finderr != nil && finderr == mongo.ErrNoDocuments {
		result, err := collection.InsertOne(context.Background(), datatypes.TypeMatrix{Matrix: matrix})
		if err != nil {
			fmt.Printf("Error saving matrix to DB: %s\n", err.Error())
			return
		}
		fmt.Printf("Typematrix created with id of %s\n", result.InsertedID)
	} else {

		filterSave := bson.M{"_id": oldmatrix.ID}

		update := bson.M{"$set": bson.M{"matrix": matrix}}

		result, err := collection.UpdateOne(context.Background(), filterSave, update)
		if err != nil {
			fmt.Printf("Error saving matrix to DB: %s\n", err.Error())
			return
		}

		fmt.Printf("Typematrix updated %d\n", result.MatchedCount)
	}

}
