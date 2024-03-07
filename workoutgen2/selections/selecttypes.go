package selections

import "math/rand"

func SelectTypes(level, minutes float32) [9]string {
	regularPos := float32(650)
	comboPos := float32(400)
	splitPos := float32(250)

	comboPos *= 1 + ((minutes - 20) / 75)
	splitPos *= 1 + ((minutes - 20) / 100)

	comboPos *= 1 + ((level - 200) / 2000)
	splitPos *= 1 + ((level - 200) / 1500)

	total := regularPos + comboPos + splitPos

	var ret [9]string

	for i := 0; i < 9; i++ {
		likelihood := rand.Float32() * total
		if likelihood < regularPos {
			ret[i] = "Regular"
		} else if likelihood < (regularPos + comboPos) {
			ret[i] = "Combo"
		} else {
			ret[i] = "Split"
		}
	}
	return ret

}
