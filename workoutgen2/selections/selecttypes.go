package selections

import "math/rand"

func SelectTypes(level, minutes float32, diff int) [9]string {
	regularPos := float32(650)
	comboPos := float32(400)
	splitPos := float32(250)

	comboPos *= 1 + 0.667*((minutes-20)/75)
	splitPos *= 1 + 0.667*((minutes-20)/100)

	comboPos *= 1 + 1.25*((level-200)/2000)
	splitPos *= 1 + 1.25*((level-200)/1500)

	if diff == 1 {
		splitPos = 0
	} else if diff == 3 {
		comboPos *= 0.875
		splitPos *= 0.875
	} else if diff == 5 || diff == 6 {
		comboPos *= 1.45
		splitPos *= 1.55
	}

	total := regularPos + comboPos + splitPos

	var ret [9]string

	for i := 0; i < 9; i++ {
		likelihood := rand.Float32() * total
		if diff == 2 || likelihood < regularPos {
			ret[i] = "Regular"
		} else if likelihood < (regularPos + comboPos) {
			ret[i] = "Combo"
		} else {
			ret[i] = "Split"
		}
	}
	return ret

}
