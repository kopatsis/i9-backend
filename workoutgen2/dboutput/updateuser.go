package dboutput

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserLast(minutes float32, difficulty int, userID string, database *mongo.Database) error {

	update := bson.M{
		"$set": bson.M{"lastmins": minutes, "lastdiff": difficulty},
	}

	var id primitive.ObjectID
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		id = oid
	} else {
		return err
	}

	filter := bson.M{"_id": id}

	collection := database.Collection("user")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
