package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateRoomHost(gid uint32, newownerpid uint32) {
	_, err := roomsCollection.UpdateOne(context.TODO(), bson.D{{"gid", gid}}, bson.D{{"$set", bson.D{{"host", int64(newownerpid)}}}})
	if err != nil {
		//panic(err)
		return
	}
}
