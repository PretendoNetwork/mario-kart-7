package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func DeletePlayerSession(rvcid uint32) {
	_, err := sessionsCollection.DeleteOne(context.TODO(), bson.D{{"rvcid", rvcid}})
	if err != nil {
		panic(err)
	}
}
