package adjustments

import "fulli9/shared"

func GeneralTyping(exerIDs [9][]string, types [9]string, exercises map[string]shared.Exercise) [3]float32 {

	genTypeToPos := map[string]int{"Legs": 0, "Core": 1, "Push": 2}

	ret := [3]float32{}
	sum := float32(0)

	for i, idList := range exerIDs {
		if types[i] == "Regular" {
			gtList := exercises[idList[0]].GeneralType
			if len(gtList) == 1 {
				ret[genTypeToPos[gtList[0]]] += 1.0
				sum += 1.0
			} else if len(gtList) == 2 {
				ret[genTypeToPos[gtList[0]]] += 2.0 / 3.0
				ret[genTypeToPos[gtList[1]]] += 2.0 / 3.0
				sum += 4.0 / 3.0
			} else {
				ret[genTypeToPos[gtList[0]]] += 1.0 / 2.0
				ret[genTypeToPos[gtList[1]]] += 1.0 / 2.0
				ret[genTypeToPos[gtList[2]]] += 1.0 / 2.0
				sum += 3.0 / 2.0
			}
		} else if types[i] == "Combo" {
			for _, id := range idList {
				gtList := exercises[id].GeneralType
				if len(gtList) == 1 {
					ret[genTypeToPos[gtList[0]]] += 1.0 / float32(len(idList))
					sum += 1.0 / float32(len(idList))
				} else if len(gtList) == 2 {
					ret[genTypeToPos[gtList[0]]] += (2.0 / 3.0) / float32(len(idList))
					ret[genTypeToPos[gtList[1]]] += (2.0 / 3.0) / float32(len(idList))
					sum += (4.0 / 3.0) / float32(len(idList))
				} else {
					ret[genTypeToPos[gtList[0]]] += (1.0 / 2.0) / float32(len(idList))
					ret[genTypeToPos[gtList[1]]] += (1.0 / 2.0) / float32(len(idList))
					ret[genTypeToPos[gtList[2]]] += (1.0 / 2.0) / float32(len(idList))
					sum += (3.0 / 2.0) / float32(len(idList))
				}
			}
		} else {
			for _, id := range idList {
				gtList := exercises[id].GeneralType
				if len(gtList) == 1 {
					ret[genTypeToPos[gtList[0]]] += 1.0 / (1.5)
					sum += 1.0 / (1.5)
				} else if len(gtList) == 2 {
					ret[genTypeToPos[gtList[0]]] += (2.0 / 3.0) / (1.5)
					ret[genTypeToPos[gtList[1]]] += (2.0 / 3.0) / (1.5)
					sum += (4.0 / 3.0) / (1.5)
				} else {
					ret[genTypeToPos[gtList[0]]] += (1.0 / 2.0) / (1.5)
					ret[genTypeToPos[gtList[1]]] += (1.0 / 2.0) / (1.5)
					ret[genTypeToPos[gtList[2]]] += (1.0 / 2.0) / (1.5)
					sum += (3.0 / 2.0) / (1.5)
				}
			}
		}
	}

	ret[0] /= sum
	ret[1] /= sum
	ret[2] /= sum

	return ret
}
