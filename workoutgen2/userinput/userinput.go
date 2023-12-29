package userinput

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetUserInputs() (float32, int, string) {
	fmt.Println("Started workout generation")
	reader := bufio.NewReader(os.Stdin)

	minutesDec := 0.00
	for minutesDec < 5 {

		fmt.Print("Time in minutes for WO as decimal (minimum 5): ")

		minutes, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading minutes input:", err)
			continue
		}

		minutes = strings.TrimSpace(minutes)
		minutesDec, err = strconv.ParseFloat(minutes, 32)
		if err != nil {
			fmt.Println("Error reading minutes input:", err)
			continue
		}

		if minutesDec < 5 {
			fmt.Println("Error reading minutes input: Not in range")
		}

	}

	difficultyInt := -1
	for difficultyInt < 0 || difficultyInt > 6 {
		fmt.Print("Difficulty from 0-6 (0=easiest, 6=hardest): ")

		difficulty, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading difficulty input:", err)
			continue
		}

		difficulty = strings.TrimSpace(difficulty)

		difficultyInt, err = strconv.Atoi(difficulty)
		if err != nil {
			fmt.Println("Error reading difficulty input:", err)
			difficultyInt = -1
			continue
		}
		if difficultyInt < 0 || difficultyInt > 6 {
			fmt.Println("Error reading difficulty input: Not in range")
		}
	}

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

	return float32(minutesDec), difficultyInt, username

}
