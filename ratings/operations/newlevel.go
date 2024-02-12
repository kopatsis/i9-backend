package operations

import (
	"fulli9/shared"
	"math"
)

func NewLevel(user shared.User, ratings [9]float32, difficulty int, count int64) (float32, float32) {
	sum := float32(0)
	for _, rating := range ratings {
		sum += rating
	}
	average := sum / 9

	newlevel := user.Level + 2 + 2*(2.5*float32(difficulty-1)-average)/float32(math.Sqrt(float64(count+2)))

	return newlevel, average
}
