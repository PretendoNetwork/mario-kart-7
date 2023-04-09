package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func AddNewUser(pid uint32) {
	_, err := usersCollection.InsertOne(context.TODO(), bson.D{{"pid", pid}, {"rating", bson.A{1000, 1000, 1000}}, {"username", GetUsernameFromPID(pid)}, {"status", "allowed"}})
	if err != nil {
		panic(err)
	}
}
