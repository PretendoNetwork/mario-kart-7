package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IsAllowedGameModeAndRegion(gamemode uint32, region uint32) bool {
	var result bson.M

	err := regionsCollection.FindOne(context.TODO(), bson.D{{Key: "id", Value: region}}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
	}

	if len(result["allowed_gamemodes"].(bson.A)) <= int(gamemode) {
		return false
	}
	if result["allowed_gamemodes"].(bson.A)[gamemode].(bool) {
		return true
	} else {
		return false
	}
}
