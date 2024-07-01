package operations

import (
	"fulli9/shared"
	"math"
)

func NewLevel(user shared.User, rating, difficulty, completedRounds int, count int64) float32 {

	baseIncrease := 2
	difficultyAdjustment := 2.5 * float32(difficulty)
	ratingAdjustment := float32(rating)
	countAdjustment := float32(math.Sqrt(float64(count + 2)))

	additionalLevel := 4 * (difficultyAdjustment - ratingAdjustment) / countAdjustment

	levelAdj := float32(baseIncrease) + additionalLevel

	if completedRounds < 9 {
		levelAdj = float32(math.Max(2, float64(levelAdj-float32(9-completedRounds))))
	}

	return user.Level + levelAdj
}
