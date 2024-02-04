package platform

import (
	"fulli9/adapts"
	"fulli9/intro"
	"fulli9/ratings"
	"fulli9/shared"
	"fulli9/userfuncs"
	"fulli9/usergeneral"
	"fulli9/views"
	"fulli9/workoutgen2"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func New() *gin.Engine {
	router := gin.Default()

	client, database, err := shared.ConnectDB()
	if err != nil {
		log.Fatalf("Error while connecting to mongoDB: %s.\nExiting.", err)
	}
	defer shared.DisConnectDB(client)

	// Won't be used
	router.GET("/", temp(database))

	// Main functionalities
	router.POST("/workouts", workoutgen2.PostWorkout(database))
	router.POST("/workouts/stretch", workoutgen2.PostStretchWorkout(database))
	router.POST("/workouts/intro", intro.PostIntroWorkout(database))
	router.POST("/workouts/rate/:id", ratings.PostRating(database))
	router.POST("/workouts/intro/rate", ratings.PostIntroRating(database))
	router.POST("/workouts/adapt/:id", adapts.PostAdaptedWorkout(database))
	router.POST("/workouts/adapt/external/:id", adapts.PostExternalAdaptedWorkout(database))

	// General User
	router.POST("/users", usergeneral.PostUser(database))
	router.PUT("/users", usergeneral.PutUser(database))
	router.GET("/users", usergeneral.GetUser(database))
	router.DELETE("/users", usergeneral.DeleteUser(database))

	// Gets/views
	router.GET("/workouts/:id", views.GetWorkout(database))
	router.GET("/workouts/external/:id", views.GetWorkout(database))
	router.GET("/workouts", views.GetWorkouts(database))
	router.GET("/stretches/:id", views.GetStrByID(database))
	router.GET("/exercises/:id", views.GetExerByID(database))
	router.GET("/stretches", views.GetStrecthes(database))
	router.GET("/exercises", views.GetExercises(database))

	// Sets/adds user specifics
	router.POST("/users/plyo", userfuncs.PostPlyo(database))
	router.POST("/users/bannedexers", userfuncs.PostBannedExer(database))
	router.POST("/users/bannedbody", userfuncs.PostBannedBody(database))
	router.POST("/users/favorites", userfuncs.PostExerFav(database))
	router.POST("/users/bannedstrs", userfuncs.PostBannedStr(database))

	// Single deletes
	router.DELETE("/users/bannedexers", userfuncs.DeleteBannedExer(database))
	router.DELETE("/users/bannedbody", userfuncs.DeleteBannedBody(database))
	router.DELETE("/users/favorites", userfuncs.DeleteExerFav(database))
	router.DELETE("/users/bannedstrs", userfuncs.DeleteBannedStr(database))

	// Clear alls
	router.DELETE("/users/bannedexers/clear", userfuncs.ClearBannedExer(database))
	router.DELETE("/users/bannedbody/clear", userfuncs.ClearBannedBody(database))
	router.DELETE("/users/favorites/clear", userfuncs.ClearExerFav(database))
	router.DELETE("/users/bannedstrs/clear", userfuncs.ClearBannedStr(database))

	return router
}

func temp(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
