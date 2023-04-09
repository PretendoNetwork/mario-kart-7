package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RemovePlayerFromRoom(gid uint32, pid uint32) {
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
		DestroyRoom(gid)
	}
}
