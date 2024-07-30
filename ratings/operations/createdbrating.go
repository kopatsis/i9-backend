package operations

import (
	"fulli9/shared"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateDatabaseRating(ratings, favorites [9]int, fullRating, fullFave int, onlyWorkout bool, workout shared.Workout) shared.StoredRating {
	ret := shared.StoredRating{
		UserID:        workout.UserID,
		WorkoutID:     workout.ID.Hex(),
		LevelAtStart:  workout.LevelAtStart,
		Minutes:       workout.Minutes,
		Difficulty:    workout.Difficulty,
		OnlyWorkout:   onlyWorkout,
		Date:          primitive.NewDateTimeFromTime(time.Now()),
		OverallRating: int(fullRating),
		OverallFave:   int(fullFave),
	}

	for i, rating := range ratings {
		if onlyWorkout {
			ret.RoundRatings[i] = shared.RoundRating{
				ActualRound: workout.Exercises[i],
				Rating:      -1,
				Fave:        -1,
				HasRating:   false,
				HasFave:     false,
			}
		}
		ret.RoundRatings[i] = shared.RoundRating{
			ActualRound: workout.Exercises[i],
			Rating:      rating,
			Fave:        fullFave,
			HasRating:   rating == -1,
			HasFave:     favorites[i] == -1,
		}
	}

	return ret
}

func CreateStrDatabaseRating(favorites []int, fullFave int, onlyWorkout bool, workout shared.StretchWorkout) shared.StoredStrRating {
	ret := shared.StoredStrRating{
		UserID:      workout.UserID,
		WorkoutID:   workout.ID.Hex(),
		Minutes:     workout.Minutes,
		Date:        primitive.NewDateTimeFromTime(time.Now()),
		OnlyWorkout: onlyWorkout,
		Fave:        fullFave,
		StaticIDs:   workout.Statics,
		DynamicIDs:  workout.Dynamics,
	}

	if !onlyWorkout {
		ret.Faves = favorites
	}

	return ret
}
