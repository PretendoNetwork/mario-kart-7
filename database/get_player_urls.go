package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayerURLs(rvcid uint32) []string {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"rvcid", rvcid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	oldurlArray := result["urls"].(bson.A)
	newurlArray := make([]string, len(oldurlArray))
	for i := 0; i < len(oldurlArray); i++ {
		fmt.Println(oldurlArray[i].(string))
		newurlArray[i] = oldurlArray[i].(string)
	}

	return newurlArray
}
