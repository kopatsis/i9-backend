package datatypes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exercise struct {
	Name         string
	Parent       string
	MinLevel     float32
	MaxLevel     float32
	PlyoRating   float32
	StartQuality float32
	BodyParts    []int
	Compatibles  map[string][2]float32
	Reps         map[int]float32
}

type Stretch struct {
	Name      string
	MinLevel  float32
	Status    string
	BodyParts []int
}

type RetExercise struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string
	Parent       string
	MinLevel     float32
	MaxLevel     float32
	PlyoRating   float32
	StartQuality float32
	BodyParts    []int
	Compatibles  map[string][2]float32
	Reps         map[int]float32
}
