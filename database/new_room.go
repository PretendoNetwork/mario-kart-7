package database

import (
	"context"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRoom(host uint32, gamemode uint32, public bool, region uint32, gameconfig uint32, playercount uint32, update uint32, communitygid uint32) uint32 {
	var gatheringId uint32
	var result bson.M

	for true {
		gatheringId = rand.Uint32() % 500000
		err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gatheringId}}, options.FindOne()).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				break
			} else {
				panic(err)
			}
		} else {
			continue
		}
	}
	players := make([]int64, 12)
	gatheringDoc := bson.D{
		{"gid", gatheringId},
		{"host", host},
		{"gamemode", gamemode},
		{"players", players},
		{"public", public},
		{"region", region},
		{"gameconfig", gameconfig},
		{"player_count", int64(0)},
		{"update", update},
		{"community_gid", communitygid}}
	_, err := roomsCollection.InsertOne(context.TODO(), gatheringDoc)
	if err != nil {
		panic(err)
	}

	return gatheringId
}
