package adjustments

import (
	"fulli9/shared"
	"time"
)

func CalcInitLevel(startLevel float32, pastWOs []shared.Workout) float32 {
	retLevel := startLevel

	lastWOTime := time.Now().AddDate(0, 0, -1)

	for _, workout := range pastWOs {
		if workout.Status != "Rated" {
			continue
		}
		if workout.Created.Time().After(lastWOTime) {
			lastWOTime = workout.Created.Time()
		}
	}

	if lastWOTime.After(time.Now().AddDate(0, 0, -1)) {
		hoursSince := float32(lastWOTime.Sub(time.Now().AddDate(0, 0, -1)).Hours())
		retLevel *= (1 - (1 / (3 * hoursSince)))
	}

	return retLevel
}

func CalcDiffLevel(diff int, level float32) float32 {
	retLevel := level
	switch diff {
	case 1:
		retLevel *= 0.7
	case 2:
		retLevel *= 0.7
	case 3:
		retLevel *= 0.875
	case 5:
		retLevel *= 1.3
	case 6:
		retLevel *= 1.45
	}
	return retLevel
}
