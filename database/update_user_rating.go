package database

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUserRating(pid uint32, gamemode uint32, rating uint32) {
	_, err := usersCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"rating." + strconv.FormatUint(uint64(gamemode), 10), rating}}}})
	if err != nil {
		panic(err)
	}
}
