package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func PositionsRequestWorkout(workout Workout, res, token string) (any, error) {
	payload := PosWorkoutRoute{
		Dynamics:     workout.Dynamics,
		Statics:      workout.Statics,
		StretchTimes: workout.StretchTimes,
		ID:           workout.ID,
		Exercises:    workout.Exercises,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	url := os.Getenv("POSITIONSURL") + "/workouts"

	if res != "" {
		url += "/" + res
	}

	return actualRequest(jsonData, url, token)
}

func PositionsRequestStrWorkout(workout StretchWorkout, res, token string) (any, error) {
	payload := PosStretchWorkoutRoute{
		Dynamics:     workout.Dynamics,
		Statics:      workout.Statics,
		StretchTimes: workout.StretchTimes,
		ID:           workout.ID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	url := os.Getenv("POSITIONSURL") + "/workouts/stretch"

	if res != "" {
		url += "/" + res
	}

	return actualRequest(jsonData, url, token)
}

func actualRequest(jsonData []byte, url, token string) (any, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var response any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
