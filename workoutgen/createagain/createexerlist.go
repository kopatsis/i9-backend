package createagain

import (
	"fulli9/workoutgen/datatypes"
	"math/rand"
	"sort"
)

func comboSelector(ratings map[string]float32, exercises map[string]datatypes.Exercise, comboID, splitID string) map[string]string {
	ret := map[string]string{"status": "COMBO"}

	// keys := make([]string, len(ratings))

	// i := 0
	// for k := range ratings {
	// 	keys[i] = k
	// 	i++
	// }

	keys := []string{}

	for k, v := range ratings {
		if exercises[k].Name != "COMBO" && exercises[k].Name != "SPLIT" && v > 0 {
			keys = append(keys, k)
		}
	}

	// randKey1 := int(rand.Float64() * float64(len(keys)))

	// for keys[randKey1] == comboID || keys[randKey1] != splitID || ratings[keys[randKey1]] == 0 {
	// 	randKey1 = int(rand.Float64() * float64(len(keys)))
	// }

	ret["id1"] = keys[int(rand.Float64()*float64(len(keys)))]

	randKey2 := int(rand.Float64() * float64(len(keys)))

	for keys[randKey2] == ret["id1"] {
		randKey2 = int(rand.Float64() * float64(len(keys)))
	}

	ret["id2"] = keys[randKey2]

	return ret
}

func splitSelector(ratings map[string]float32, exercises map[string]datatypes.Exercise, comboID, splitID string) map[string]string {
	ret := map[string]string{"status": "SPLIT"}

	keys := []string{}

	for k, v := range ratings {
		if exercises[k].Name != "COMBO" && exercises[k].Name != "SPLIT" && exercises[k].Compatibles != nil && len(exercises[k].Compatibles) > 0 && v > 0 {
			keys = append(keys, k)
		}
	}

	ret["id1"] = keys[int(rand.Float64()*float64(len(keys)))]

	keys2 := make([]string, len(exercises[ret["id1"]].Compatibles))

	i := 0
	for k := range exercises[ret["id1"]].Compatibles {
		keys2[i] = k
		i++
	}

	ret["id2"] = keys2[int(rand.Float64()*float64(len(keys2)))]

	return ret
}

func CreateExerList(exercises map[string]datatypes.Exercise, ratings map[string]float32) [9]map[string]string {
	var ret [9]map[string]string

	var comboID string
	var splitID string
	for _, exercise := range exercises {
		if exercise.Name == "COMBO" {
			comboID = exercise.ID.Hex()
		} else if exercise.Name == "SPLIT" {
			splitID = exercise.ID.Hex()
		}
	}

	var total float64
	for _, rating := range ratings {
		total += float64(rating)
	}

	var randInds []float64
	for i := 0; i < 9; i++ {
		randInds = append(randInds, rand.Float64()*total)
	}

	sort.Float64s(randInds)

	index := 0
	for id, rating := range ratings {
		total -= float64(rating)
		for index < 9 && randInds[index] <= total {
			if id == comboID {
				ret[index] = comboSelector(ratings, exercises, comboID, splitID)
			} else if id == splitID {
				ret[index] = splitSelector(ratings, exercises, comboID, splitID)
			} else {
				ret[index] = map[string]string{"id": id, "Status": "Normal"}
			}
		}
	}

	return ret
}
