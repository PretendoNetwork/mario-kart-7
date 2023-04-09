package database

import (
	"context"
	"encoding/base64"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCommunity(host uint32, communityType uint32, password string, attribs []uint32, applicationBuffer []byte, participationStartDate uint64, participationEndDate uint64) uint32 {
	var gatheringId uint32
	var result bson.M

	for true {
		gatheringId = rand.Uint32() % 500000
		err := communitiesCollection.FindOne(context.TODO(), bson.D{{"gid", gatheringId}}, options.FindOne()).Decode(&result)
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
	participants := make([]int64, 1)
	participants[0] = int64(host)
	gatheringDoc := bson.D{
		{"gid", gatheringId},
		{"host", host},
		{"type", communityType},
		{"password", password},
		{"attribs", attribs},
		{"application_buffer", base64.StdEncoding.EncodeToString(applicationBuffer)},
		{"start_date", participationStartDate},
		{"end_date", participationEndDate},
		{"sessions", 0},
		{"participants", participants}}
	_, err := communitiesCollection.InsertOne(context.TODO(), gatheringDoc)
	if err != nil {
		panic(err)
	}

	return gatheringId
}
