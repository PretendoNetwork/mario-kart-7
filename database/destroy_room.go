package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func DestroyRoom(gid uint32) {
	_, err := roomsCollection.DeleteOne(context.TODO(), bson.D{{"gid", gid}})
	if err != nil {
		//panic(err)
		return
	}
}
