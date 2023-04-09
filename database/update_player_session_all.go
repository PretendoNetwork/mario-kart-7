package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePlayerSessionAll(rvcid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"rvcid", rvcid}}, bson.D{{"$set", bson.D{{"rvcid", rvcid}, {"urls", urls}, {"ip", ip}, {"port", port}}}})
	if err != nil {
		panic(err)
	}
}
