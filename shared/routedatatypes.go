package shared

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkoutRoute struct {
	Time       float32 `json:"time" binding:"required,min=10,max=240"`
	Difficulty int     `json:"diff" binding:"required,min=1,max=6"`
}

type StrWorkoutRoute struct {
	Time float32 `json:"time" binding:"required"`
}

type PatchWorkout struct {
	PausedMinutes float32 `json:"minutes" binding:"required"`
	Status        string  `json:"status" binding:"required"`
}

type IntroWorkoutRoute struct {
	Time float32 `json:"time" binding:"required,min=25,max=60"`
}

type AdaptWorkoutRoute struct {
	Difficulty int  `json:"difficulty" binding:"required,min=1,max=6"`
	AsNew      bool `json:"asnew"`
}

type RateIntroRoute struct {
	Rounds float32 `json:"rounds"`
}

type RateRoute struct {
	Ratings    []float32 `json:"ratings" binding:"required"`
	Favoritism []float32 `json:"faves" binding:"required"`
}

type UserRoute struct {
	Name string `json:"name"`
}
type MergeRoute struct {
	LocalJWT string `json:"localjwt"`
}

type PlyoRoute struct {
	Plyo int `json:"plyo" binding:"required,min=0,max=5"`
}

type PushupSettingRoute struct {
	Setting string `json:"pushupsetting" binding:"required"`
}

type PaySettingRoute struct {
	Paying bool `json:"paying" binding:"required"`
}

type ExerListRoute struct {
	ExerList []string `json:"exerlist" binding:"required"`
}

type StrListRoute struct {
	StrList []string `json:"strlist" binding:"required"`
}

type BodyListRoute struct {
	BodyList []int `json:"bodylist" binding:"required"`
}

type ExerMapRoute struct {
	ExerMap map[string]float32 `json:"exermap" binding:"required"`
}

type PosStretchWorkoutRoute struct {
	Dynamics     []string
	Statics      []string
	StretchTimes StretchTimes
	ID           primitive.ObjectID
}

type PosWorkoutRoute struct {
	Dynamics     []string
	Statics      []string
	StretchTimes StretchTimes
	ID           primitive.ObjectID
	Exercises    [9]WorkoutRound
}
