package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func AddPlayerSession(rvcid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.InsertOne(context.TODO(), bson.D{{"rvcid", rvcid}, {"urls", urls}, {"ip", ip}, {"port", port}})
	if err != nil {
		panic(err)
	}
}
