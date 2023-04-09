package database

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindRoom(gamemode uint32, public bool, region uint32, gameconfig uint32, vacantcount uint32, update uint32, communitygid uint32) uint32 {
	var result bson.M
	maxplayersinroom := 12 - vacantcount
	filter := bson.D{
		{"gamemode", gamemode},
		{"public", public},
		{"region", region},
		{"gameconfig", gameconfig},
		{"player_count", bson.D{{"$lte", maxplayersinroom}}},
		{"update", update},
		{"community_gid", communitygid}}

	err := roomsCollection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return math.MaxUint32
		} else {
			panic(err)
		}
	} else {
		return uint32(result["gid"].(int64))
	}
}
