package alteredfuncs

import "math/rand"

func SelectTypes(levelSteps []float32, minutes float32) [9]string {
	regularPos := float32(600)
	comboPos := float32(200)
	splitPos := float32(200)

	comboPos *= 1 + ((minutes - 20) / 75)
	splitPos *= 1 + ((minutes - 20) / 100)

	var ret [9]string

	for i, level := range levelSteps {

		if i < 3 {
			ret[i] = "Regular"
		} else {
			tempCombo := (1 + ((level - 200) / 2000)) * comboPos
			tempSplit := (1 + ((level - 200) / 1500)) * splitPos

			total := regularPos + tempCombo + tempSplit

			likelihood := rand.Float32() * total
			if likelihood < regularPos {
				ret[i] = "Regular"
			} else if likelihood < (regularPos + tempCombo) {
				ret[i] = "Combo"
			} else {
				ret[i] = "Split"
			}
		}

	}

	return ret
}
