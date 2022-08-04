package main

import (
	"context"
	//"math"
	"fmt"
	"math/rand"

	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newCommunity(host uint32, communityType uint32, password string, attribs []uint32, applicationBuffer []byte, participationStartDate uint64, participationEndDate uint64) uint32 {
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

func communityExists(gid uint32) (bool) {
	var result bson.M

	err := communitiesCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}

	return false
}

func getCommunityInfo(gid uint32) (uint32, uint32, string, []uint32, []byte, uint64, uint64, uint32, uint32) {
	var result bson.M

	err := communitiesCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}
	attribs := make([]uint32, len(result["attribs"].(bson.A)))
	for index, attrib := range result["attribs"].(bson.A) {
		if val, ok := attrib.(uint32); ok {
			attribs[index] = val
		}
	}
	
	application_buffer, _ := base64.StdEncoding.DecodeString(result["application_buffer"].(string))

	return uint32(result["host"].(int64)), uint32(result["type"].(int64)), result["password"].(string), attribs, application_buffer, uint64(result["start_date"].(int64)), uint64(result["end_date"].(int64)), uint32(result["sessions"].(int32)), uint32(len(result["participants"].(bson.A)))
}

func findCommunitiesWithParticipant(pid uint32) []uint32 {
	arr := []uint32{pid}
	gatheringIDs := []uint32{}

	cur, err := communitiesCollection.Find(context.TODO(), bson.M{"participants": bson.M{"$in": arr}}, options.Find())
	if err != nil {
		return nil
	}

    for cur.Next(context.TODO()) {
		fmt.Println("test")
        //Create a value into which the single document can be decoded
        var result bson.M
        err := cur.Decode(&result)
        if err != nil {
			return nil
        }

        gatheringIDs = append(gatheringIDs, (uint32)(result["gid"].(int64)))

    }

	return gatheringIDs
}
