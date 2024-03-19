package operations

import "fulli9/shared"

func AdjustRatings(ratings [9]float32, workout shared.Workout) [9]float32 {

	ret := [9]float32{}

	cardioToRatingExpMap := map[int]float32{0: 2.5, 1: 3.125, 2: 3.75, 3: 4.5, 4: 5.25, 5: 6, 6: 6.625, 7: 7.25, 8: 7.75, 9: 8.5}

	for i, rating := range ratings {
		expected := cardioToRatingExpMap[int(workout.CardioRatings[i])]

		realRating := rating*(1.0/3.0) + (rating-expected+5)*(2.0/3.0)

		ret[i] = realRating
	}

	return ret
}
