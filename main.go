package main

import (
	"bufio"
	"fmt"
	"fulli9/esa"
	"fulli9/intro"
	"fulli9/ratings"
	"fulli9/views"
	"fulli9/workoutgen2"
	"fulli9/workoutgen2/userinput"
	"os"
	"strconv"
	"strings"
)

func main() {
	displayMenu()
	userInput := getUserInput()
	for userInput == -1 {
		userInput = getUserInput()
	}

	switch userInput {
	case 1:
		minutes, diff, username := userinput.GetUserInputs()
		if WO, err := workoutgen2.WorkoutGen(minutes, diff, username); err != nil {
			fmt.Println(err)
		} else {
			views.PrettyPrint(WO)
		}
	case 2:
		username := getUserName()
		if WO, err := intro.GenerateIntroWorkout(username); err != nil {
			fmt.Println(err)
		} else {
			views.PrettyPrint(WO)
		}
	case 3:
		username := getUserName()
		if err := ratings.RateWorkout(username); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Success rating WO")
		}
	case 4:
		username := getUserName()
		rounds := getRounds(username)
		if err := ratings.RateIntroWorkout(username, rounds); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Success rating Intro WO")
		}
	case 5:
		esa.RunESA()
	case 6:
		username := getUserName()
		views.ViewLastWorkout(username)
	default:
		fmt.Println("Error in logic for menu options. Restart")
	}
}

func displayMenu() {
	fmt.Println("Options for i9 basic:")
	fmt.Println("1. Generate WO (reqs username, difficulty, time)")
	fmt.Println("2. Intro WO (reqs username)")
	fmt.Println("3. Rate last WO (reqs username)")
	fmt.Println("4. Rate Intro WO (reqs username, rounds completed)")
	fmt.Println("5. Run esa (no reqs)")
	fmt.Println("6. View last WO for user (reqs username)")
}

func getUserInput() int {
	fmt.Print("Enter your choice (1-6): ")
	var userInput string
	fmt.Scanln(&userInput)
	if num, err := strconv.Atoi(userInput); err == nil {
		if num >= 1 && num <= 5 {
			return num
		}
	}
	return -1
}

func getUserName() string {
	reader := bufio.NewReader(os.Stdin)
	username := ""
	for username == "" {
		fmt.Print("Username: ")
		usernameTMP, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading username input:", err)
			continue
		}

		username = strings.TrimSpace(usernameTMP)
		if username == "" {
			fmt.Println("Error reading username input: Can't be blank")
		}
	}
	return username
}

func getRounds(username string) float32 {
	reader := bufio.NewReader(os.Stdin)
	rounds := float32(-1)
	for rounds < 0 {
		fmt.Printf("Number of rounds (0.0-9.0 as float) completed by %s: ", username)
		roundsSt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading rounds complete input:", err)
			continue
		}

		roundsTMP, err := strconv.ParseFloat(roundsSt, 32)
		if err != nil {
			fmt.Println("Error reading rounds complete input: Must be float32")
			continue
		}

		roundsFL32 := float32(roundsTMP)
		if roundsFL32 < 0 || roundsFL32 > 9 {
			fmt.Println("Error reading rounds complete input: Incorrect range")
			continue
		} else {
			rounds = roundsFL32
		}

	}
	return rounds
}
