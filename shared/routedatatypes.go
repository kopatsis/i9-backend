package shared

type UserIDRoute struct {
	UserID string `json:"userid"`
}

type WorkoutRoute struct {
	UserID     string  `json:"userid"`
	Time       float32 `json:"time" binding:"required"`
	Difficulty int     `json:"diff" binding:"required"`
}

type StrWorkoutRoute struct {
	UserID string  `json:"userid"`
	Time   float32 `json:"time" binding:"required"`
}

type IntroWorkoutRoute struct {
	UserID string  `json:"userid"`
	Time   float32 `json:"time" binding:"required"`
}

type AdaptWorkoutRoute struct {
	UserID     string `json:"userid"`
	Difficulty int    `json:"difficulty" binding:"required"`
}

type RateIntroRoute struct {
	UserID string  `json:"userid"`
	Rounds float32 `json:"rounds" binding:"required"`
}

type RateRoute struct {
	UserID  string    `json:"userid"`
	Ratings []float32 `json:"ratings" binding:"required"`
}

type PostUserRoute struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"username" binding:"required"`
}

type PutUserRoute struct {
	UserID   string `json:"userid"`
	Name     string `json:"name"`
	UserName string `json:"username"`
}

type PlyoRoute struct {
	UserID string `json:"userid"`
	Plyo   int    `json:"plyo" binding:"required"`
}

type ExerListRoute struct {
	UserID   string   `json:"userid"`
	ExerList []string `json:"exerlist" binding:"required"`
}

type StrListRoute struct {
	UserID  string   `json:"userid"`
	StrList []string `json:"strlist" binding:"required"`
}

type BodyListRoute struct {
	UserID   string `json:"userid"`
	BodyList []int  `json:"bodylist" binding:"required"`
}

type ExerMapRoute struct {
	UserID  string             `json:"userid"`
	ExerMap map[string]float32 `json:"exermap" binding:"required"`
}
