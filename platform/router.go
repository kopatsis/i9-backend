package platform

import (
	"context"
	"fulli9/adapts"
	"fulli9/intro"
	"fulli9/platform/cache"
	"fulli9/platform/middleware"
	"fulli9/ratings"
	"fulli9/shared"
	"fulli9/userfuncs"
	"fulli9/usergeneral"
	"fulli9/views"
	"fulli9/workoutgen2"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(database *mongo.Database, firebase *firebase.App, boltDB *bbolt.DB) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.JWTAuthMiddleware(firebase))

	// Won't be used
	router.GET("/", temp())
	router.GET("/tpv/:id", tpv(database)) // Temp for personal viewing a workout by id

	// Main functionalities
	router.POST("/workouts", workoutgen2.PostWorkout(database, boltDB))
	router.POST("/workouts/stretch", workoutgen2.PostStretchWorkout(database, boltDB))
	router.POST("/workouts/intro", intro.PostIntroWorkout(database, boltDB))

	// Rating functionalities
	router.POST("/workouts/rate/:id", ratings.PostRating(database, boltDB))
	router.POST("/workouts/intro/rate", ratings.PostIntroRating(database))
	router.POST("/users/quiz", ratings.PostQuiz(database))

	// Restart (ptl adapt) functionalities
	router.POST("/workouts/restart/:id", adapts.PostRestartedWorkout(database, boltDB))
	router.POST("/workouts/stretch/restart/:id", adapts.PostRestartedStrWorkout(database))
	router.POST("/workouts/intro/restart/:id", adapts.PostRestartedIntroWorkout(database))

	// Mildly deprecated functionalities
	router.POST("/workouts/adapt/:id", adapts.PostAdaptedWorkout(database, boltDB))
	router.POST("/workouts/adapt/external/:id", adapts.PostExternalAdaptedWorkout(database, boltDB))
	router.POST("/workouts/clone/:id", adapts.CloneWorkoutHandler(database))
	router.POST("/workouts/stretch/clone/:id", adapts.CloneStretchWorkoutHandler(database))

	// Discard and retries
	router.POST("/workouts/retry/:id", workoutgen2.PostWorkoutRetry(database, boltDB))
	router.POST("/workouts/intro/retry/:id", intro.PostIntroWorkoutRetry(database, boltDB))
	router.POST("/workouts/stretch/retry/:id", workoutgen2.PostStretchWorkoutRetry(database, boltDB))

	// Patching workout
	router.PATCH("/workouts/:id", workoutgen2.PatchWorkout(database))
	router.PATCH("/workouts/stretch/:id", workoutgen2.PatchStretchWorkout(database))
	router.PATCH("/rename/:id", workoutgen2.Rename(database))
	router.PATCH("/workouts/pin/:id", workoutgen2.PinWorkout(database))
	router.PATCH("/workouts/stretch/pin/:id", workoutgen2.PinStretchWorkout(database))

	// General User
	router.POST("/users/local", usergeneral.PostLocalUser(database))
	router.POST("/users", usergeneral.PostUser(database))
	router.PATCH("/users", usergeneral.PatchUser(database))
	router.PATCH("/users/merge", usergeneral.MergeLocalUser(database))
	router.GET("/users", usergeneral.GetUser(database))
	router.GET("/users/local/:id", usergeneral.GetLocalJWT(database))
	router.GET("/refresh", usergeneral.GetToken(database))
	router.DELETE("/users", usergeneral.DeleteUser(database))

	// Gets/views
	router.GET("/workouts/:id", views.GetWorkout(database))
	router.GET("/workouts/external/:id", views.GetWorkout(database))
	router.GET("/workouts", views.GetWorkouts(database))
	router.GET("/workouts/stretch/:id", views.GetStretchWorkout(database))
	router.GET("/workouts/stretch", views.GetStretchWorkouts(database))
	router.GET("/history", views.GetHistory(database))
	router.GET("/recent", views.GetMostRecent(database))

	// Gets/views 2
	router.GET("/stretches/:id", views.GetStrByID(database, boltDB))
	router.GET("/exercises/:id", views.GetExerByID(database, boltDB))
	router.GET("/stretches", views.GetStrecthes(database))
	router.GET("/exercises", views.GetExercises(database))
	router.GET("/library", views.GetLibrary(database, boltDB))

	// Sets user specifics
	router.PATCH("/users/pushup", userfuncs.PatchPushupSetting(database))
	router.PATCH("/users/plyo", userfuncs.PatchPlyo(database))
	router.PATCH("/users/paying", userfuncs.PatchPlyo(database))

	// Adds user specifics
	router.PATCH("/users/bannedexers", userfuncs.PatchBannedExer(database))
	router.PATCH("/users/bannedbody", userfuncs.PatchBannedBody(database))
	router.PATCH("/users/favorites", userfuncs.PatchExerFav(database))
	router.PATCH("/users/bannedstrs", userfuncs.PatchBannedStr(database))

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

	// Clears cache
	router.DELETE("/clearcache", cache.ClearCache(boltDB))

	return router
}

func temp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}

func tpv(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var workout shared.Workout

		idStr, exists := c.Params.Get("id")
		if !exists {
			c.JSON(400, gin.H{
				"Error": "Issue with param",
				"Exact": "Unable to get ID from URL parameter",
			})
			return
		}

		var id primitive.ObjectID
		if oid, err := primitive.ObjectIDFromHex(idStr); err == nil {
			id = oid
		} else {
			c.JSON(400, gin.H{
				"Error": "Issue with user ID",
				"Exact": err.Error(),
			})
			return
		}

		collection := database.Collection("workout")
		filter := bson.D{{Key: "_id", Value: id}}

		err := collection.FindOne(context.Background(), filter).Decode(&workout)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing workout",
				"Exact": err.Error(),
			})
			return
		}

		var filterEx bson.M

		collection = database.Collection("exercise")

		cursor, err := collection.Find(context.Background(), filterEx)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exercises",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastExers []shared.Exercise
		err = cursor.All(context.Background(), &pastExers)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exercises",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastExers) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing exercises",
				"Exact": "No results returned (check type)",
			})
			return
		}

		var filterStr bson.M

		collection = database.Collection("stretch")

		cursor, err = collection.Find(context.Background(), filterStr)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": err.Error(),
			})
			return
		}
		defer cursor.Close(context.Background())

		var pastStrs []shared.Stretch
		err = cursor.All(context.Background(), &pastStrs)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": err.Error(),
			})
			return
		}

		if len(pastStrs) == 0 {
			c.JSON(400, gin.H{
				"Error": "Issue with viewing stretch",
				"Exact": "No results returned (check type)",
			})
			return
		}

		namedStatics := []string{}
		for _, staticID := range workout.Statics {
			name := ""
			for _, str := range pastStrs {
				if str.ID.Hex() == staticID {
					name = str.Name
					break
				}
			}
			namedStatics = append(namedStatics, name)
		}

		workout.Statics = namedStatics

		namedDynamics := []string{}
		for _, dynamicID := range workout.Dynamics {
			name := ""
			for _, str := range pastStrs {
				if str.ID.Hex() == dynamicID {
					name = str.Name
					break
				}
			}
			namedDynamics = append(namedDynamics, name)
		}

		workout.Dynamics = namedDynamics

		for i, round := range workout.Exercises {
			namelist := []string{}
			for _, exer := range round.ExerciseIDs {
				name := ""
				for _, ex := range pastExers {
					if ex.ID.Hex() == exer {
						name = ex.Name
						break
					}
				}
				namelist = append(namelist, name)
			}
			workout.Exercises[i].ExerciseIDs = namelist
		}

		c.JSON(http.StatusOK, workout)
	}
}
