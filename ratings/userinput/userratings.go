package userinput

import (
	"bufio"
	"fmt"
	"fulli9/shared"
	"os"
	"strconv"
	"strings"
)

func GetUserRatings(workout shared.Workout, allexer map[string]shared.Exercise) [9]float32 {

	fmt.Printf("Rating process for workout: ID: %s, Name: %s\n\n", workout.ID, workout.Name)
	reader := bufio.NewReader(os.Stdin)

	var ret [9]float32

	for i := 0; i < 9; i++ {

		nameList := ""
		for _, id := range workout.Exercises[i].ExerciseIDs {
			nameList += allexer[id].Name
			nameList += ", "
		}

		fmt.Printf("Rate round %d: %s %s\n", i+1, nameList, workout.Exercises[i].Status)

		currentRating := float32(-1)

		for currentRating < 0 || currentRating > 10 {
			fmt.Print("Rating 0-10 as decimal (0=easy, 10=hard): ")

			rating, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading rating input:", err)
				continue
			}

			rating = strings.TrimSpace(rating)
			ratingDec, err := strconv.ParseFloat(rating, 32)
			if err != nil {
				fmt.Println("Error reading rating input:", err)
				continue
			}

			currentRating = float32(ratingDec)
		}

		ret[i] = currentRating
	}

	return ret

}
