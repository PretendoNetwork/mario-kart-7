package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayerSessionAddress(rvcid uint32) string {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"rvcid", rvcid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(rvcid)
		//panic(err)
	}

	return result["ip"].(string) + ":" + result["port"].(string)
}
