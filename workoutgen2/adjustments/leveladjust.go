package adjustments

import (
	"fulli9/workoutgen2/datatypes"
	"time"
)

func CalcNewLevel(difficulty int, startLevel float32, pastWOs []datatypes.Workout) float32 {
	retLevel := startLevel
	switch difficulty {
	case 1:
		retLevel *= 0.85
	case 2:
		retLevel *= 0.925
	case 4:
		retLevel *= 1.075
	case 5:
		retLevel *= 1.15
	case 6:
		retLevel *= 1.3
	default:
		return retLevel
	}

	lastWOTime := time.Now().AddDate(0, 0, -1)

	for _, workout := range pastWOs {
		if workout.Date.Time().After(lastWOTime) {
			lastWOTime = workout.Date.Time()
		}
	}

	if lastWOTime.After(time.Now().AddDate(0, 0, -1)) {
		hoursSince := float32(lastWOTime.Sub(time.Now().AddDate(0, 0, -1)).Hours())
		retLevel *= (1 - (1 / (3 * hoursSince)))
	}

	return retLevel
}
