package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRoomPlayers(gid uint32) []uint32 {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	dbPlayerList := result["players"].(bson.A)
	pidList := make([]uint32, 0)

	for i := 0; i < 12; i++ {
		if (uint32)(dbPlayerList[i].(int64)) != 0 {
			pidList = append(pidList, (uint32)(dbPlayerList[i].(int64)))
		}
	}

	return pidList
}
