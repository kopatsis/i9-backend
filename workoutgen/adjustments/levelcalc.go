package adjustments

import (
	"fulli9/workoutgen/datatypes"
	"time"
)

func CalcNewLevel(difficulty int, startLevel float64, pastWOs []datatypes.Workout) float64 {
	retLevel := startLevel
	switch difficulty {
	case 1:
		retLevel *= 0.75
	case 2:
		retLevel *= 0.875
	case 4:
		retLevel *= 1.15
	case 5:
		retLevel *= 1.3
	case 6:
		retLevel *= 1.5
	default:
		return retLevel
	}

	lastWOTime := time.Now().AddDate(0, 0, -1)

	for _, workout := range pastWOs {
		if workout.Created.Time().After(lastWOTime) {
			lastWOTime = workout.Created.Time()
		}
	}

	if lastWOTime.After(time.Now().AddDate(0, 0, -1)) {
		hoursSince := lastWOTime.Sub(time.Now().AddDate(0, 0, -1)).Hours()
		retLevel *= (1 - (1 / (3 * hoursSince)))
	}

	return retLevel
}
