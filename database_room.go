package main

import (
	"context"
	"math"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newRoom(host uint32, gamemode uint32, public bool, region uint32, gameconfig uint32, playercount uint32, update uint32, communitygid uint32) uint32 {
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

func addPlayerToRoom(gid uint32, pid uint32, addplayercount uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		//panic(err)
		return
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
		//panic(err)
		return
	}
}

func removePlayerFromRoom(gid uint32, pid uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		//panic(err)
		return
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
		//panic(err)
		return
	}

	if newplayercount <= 0 {
		destroyRoom(gid)
	}
}

func removePlayer(pid uint32) {
	var result bson.M
	arr := []uint32{pid}

	err := roomsCollection.FindOne(context.TODO(), bson.M{"players": bson.M{"$in": arr}}, options.FindOne()).Decode(&result)
	if err != nil {
		return
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

	_, err = roomsCollection.UpdateOne(context.TODO(), bson.D{{"gid", result["gid"]}}, bson.D{{"$set", bson.D{{"players", newPlayerList}, {"player_count", newplayercount}}}})
	if err != nil {
		return
		//panic(err)
	}

	if newplayercount <= 0 {
		destroyRoom((uint32)(result["gid"].(int64)))
	}
}

func destroyRoom(gid uint32) {
	_, err := roomsCollection.DeleteOne(context.TODO(), bson.D{{"gid", gid}})
	if err != nil {
		//panic(err)
		return
	}
}

func findRoom(gamemode uint32, public bool, region uint32, gameconfig uint32, vacantcount uint32, update uint32, communitygid uint32) uint32 {
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

func getRoomInfo(gid uint32) (uint32, uint32, uint32, uint32, uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	return uint32(result["host"].(int64)), uint32(result["gamemode"].(int64)), uint32(result["region"].(int64)), uint32(result["gameconfig"].(int64)), uint32(result["update"].(int64))
}

func getRoomPlayers(gid uint32) ([]uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	dbPlayerList := result["players"].(bson.A)
	pidList := make([]uint32, 0)

	for i := 0; i < 12; i++ {
		if((uint32)(dbPlayerList[i].(int64)) != 0){
			pidList = append(pidList, (uint32)(dbPlayerList[i].(int64)))
		}
	}

	return pidList
}

func updateRoomHost(gid uint32, newownerpid uint32) {
	_, err := roomsCollection.UpdateOne(context.TODO(), bson.D{{"gid", gid}}, bson.D{{"$set", bson.D{{"host", int64(newownerpid)}}}})
	if err != nil {
		//panic(err)
		return
	}
}
