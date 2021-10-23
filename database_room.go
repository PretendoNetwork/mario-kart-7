package main

import (
	"context"
	"math"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newRoom(host uint32, gamemode uint32, public bool, region uint32, gameconfig uint32, playercount uint32, dlcmode uint32) uint32 {
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
		{"dlc_mode", dlcmode}}
	_, err := roomsCollection.InsertOne(context.TODO(), gatheringDoc)
	if err != nil {
		panic(err)
	}

	return gatheringId
}

func addPlayerToRoom(gid uint32, pid uint32, addplayercount uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	oldPlayerList := result["players"].(bson.A)
	newPlayerList := make([]int64, 12)
	for i := 0; i < 12; i++ {
		if oldPlayerList[i].(int64) == int64(pid) || oldPlayerList[i].(int64) == -1*int64(pid) {
			newPlayerList[i] = 0
		} else {
			newPlayerList[i] = oldPlayerList[i].(int64)
		}
	}
	unassignedPlayers := addplayercount
	needToAddGuest := (unassignedPlayers > 1)
	for i := 0; i < 12; i++ {
		if newPlayerList[i] == 0 && unassignedPlayers > 0 {
			if unassignedPlayers == 1 && needToAddGuest {
				newPlayerList[i] = -1 * int64(pid)
				needToAddGuest = false
			} else {
				newPlayerList[i] = int64(pid)
			}
			unassignedPlayers--
		}
	}

	newplayercount := result["player_count"].(int64) + int64(addplayercount)

	_, err = roomsCollection.UpdateOne(context.TODO(), bson.D{{"gid", gid}}, bson.D{{"$set", bson.D{{"players", newPlayerList}, {"player_count", newplayercount}}}})
	if err != nil {
		panic(err)
	}
}

func removePlayerFromRoom(gid uint32, pid uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	oldPlayerList := result["players"].(bson.A)
	newPlayerList := make([]int64, 12)
	newplayercount := result["player_count"].(int64)
	for i := 0; i < 12; i++ {
		newPlayerList[i] = oldPlayerList[i].(int64)
		if newPlayerList[i] == int64(pid) || newPlayerList[i] == -1*int64(pid) {
			newPlayerList[i] = 0
			newplayercount--
		}
	}

	_, err = roomsCollection.UpdateOne(context.TODO(), bson.D{{"gid", gid}}, bson.D{{"$set", bson.D{{"players", newPlayerList}, {"player_count", newplayercount}}}})
	if err != nil {
		panic(err)
	}

	if newplayercount <= 0 {
		destroyRoom(gid)
	}
}

func destroyRoom(gid uint32) {
	_, err := roomsCollection.DeleteOne(context.TODO(), bson.D{{"gid", gid}})
	if err != nil {
		panic(err)
	}
}

func findRoom(gamemode uint32, public bool, region uint32, gameconfig uint32, vacantcount uint32, dlcmode uint32) uint32 {
	var result bson.M
	maxplayersinroom := 12 - vacantcount
	filter := bson.D{
		{"gamemode", gamemode},
		{"public", public},
		{"region", region},
		{"gameconfig", gameconfig},
		{"player_count", bson.D{{"$lte", maxplayersinroom}}},
		{"dlc_mode", dlcmode}}

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

func getRoomInfo(gid uint32) (uint32, uint32, uint32, uint32, uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	return uint32(result["host"].(int64)), uint32(result["gamemode"].(int64)), uint32(result["region"].(int64)), uint32(result["gameconfig"].(int64)), uint32(result["dlc_mode"].(int64))
}

func updateRoomHost(gid uint32, newownerpid uint32) {
	_, err := roomsCollection.UpdateOne(context.TODO(), bson.D{{"gid", gid}}, bson.D{{"$set", bson.D{{"host", int64(newownerpid)}}}})
	if err != nil {
		panic(err)
	}
}
