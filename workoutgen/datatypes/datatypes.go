package datatypes

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	Name                  string             `bson:"name"`
	Username              string             `bson:"username"`
	Level                 float64            `bson:"level"`
	BannedExercises       []string           `bson:"bannedExer"`
	BannedParts           []int              `bson:"bannedParts"`
	PlyoTolerance         int                `bson:"plyoToler"`
	FavoriteExercises     []string           `bson:"favoriteExer"`
	ExerSpecModifications map[string]float64 `bson:"exerMods"`
}

type Workout struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	UserID       string             `bson:"userid"`
	Created      primitive.DateTime `bson:"created"`
	Status       string             `bson:"status"`
	Times        map[string]float32 `bson:"times"`
	LevelAtStart float64            `bson:"level"`
	Dynamics     []string           `bson:"dynamics"`
	Statics      []string           `bson:"statics"`
	Exercises    []map[string]any   `bson:"exercises"`
}

// type WorkoutPract struct {
// 	ID       primitive.ObjectID `bson:"_id,omitempty"`
// 	Name     string             `bson:"name"`
// 	Username string             `bson:"username"`
// 	Created  primitive.DateTime `bson:"created"`
// }

type Exercise struct {
	ID           primitive.ObjectID    `bson:"_id,omitempty"`
	Name         string                `bson:"name"`
	Parent       string                `bson:"parent"`
	MinLevel     float32               `bson:"minlevel"`
	MaxLevel     float32               `bson:"maxlevel"`
	PlyoRating   float32               `bson:"plyorating"`
	StartQuality float32               `bson:"startquality"`
	BodyParts    []int                 `bson:"bodyparts"`
	Compatibles  map[string][2]float32 `bson:"compatibles"`
	Reps         map[int]float32       `bson:"reps"`
}

type Stretch struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	MinLevel  float32            `bson:"minlevel"`
	Status    string             `bson:"status"`
	BodyParts []int              `bson:"bodyparts"`
}
