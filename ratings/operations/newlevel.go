package operations

import (
	"fulli9/workoutgen2/datatypes"
	"math"
)

func NewLevel(user datatypes.User, ratings [9]float32, difficulty int, count int64) float32 {
	sum := float32(0)
	for _, rating := range ratings {
		sum += rating
	}
	average := sum / 9

	newlevel := user.Level + 2 + 4*(2*float32(difficulty-1)-average)/float32(math.Sqrt(float64(count+2)))

	return newlevel
}
