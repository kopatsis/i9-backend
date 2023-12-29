package stretches

import (
	"math"
	"math/rand"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetStretchWO(level, minutes float64, database *mongo.Database) (map[string][]string, map[string]float64) {
	ret := map[string][]string{}

	allStatic, allDynamic := StretchesFromDB(database, level)

	stretchSecs := (60 * minutes) / 2

	stretchSets := math.Round(stretchSecs / 20)

	stretchSecsPerSet := stretchSecs / stretchSets

	ret["static"] = []string{}
	ret["dynamic"] = []string{}

	for i := 0; i < int(stretchSets); i++ {
		ret["static"] = append(ret["static"], allStatic[int(rand.Float64()*float64(len(allStatic)))].ID.Hex())
		ret["dynamic"] = append(ret["static"], allStatic[int(rand.Float64()*float64(len(allDynamic)))].ID.Hex())
	}

	times := map[string]float64{}
	times["secs"] = stretchSecs
	times["sets"] = stretchSets
	times["secsperset"] = stretchSecsPerSet

	return ret, times
}
