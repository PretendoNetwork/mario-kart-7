package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRoomInfo(gid uint32) (uint32, uint32, uint32, uint32, uint32) {
	var result bson.M

	err := roomsCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	return uint32(result["host"].(int64)), uint32(result["gamemode"].(int64)), uint32(result["region"].(int64)), uint32(result["gameconfig"].(int64)), uint32(result["update"].(int64))
}
