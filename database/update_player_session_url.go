package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdatePlayerSessionURL(rvcid uint32, oldurl string, newurl string) {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"rvcid", rvcid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	oldurlArray := result["urls"].(bson.A)
	newurlArray := make([]string, len(oldurlArray))
	for i := 0; i < len(oldurlArray); i++ {
		if oldurlArray[i].(string) == oldurl {
			newurlArray[i] = newurl
		} else {
			newurlArray[i] = oldurlArray[i].(string)
		}
	}

	_, err = sessionsCollection.UpdateOne(context.TODO(), bson.D{{"rvcid", rvcid}}, bson.D{{"$set", bson.D{{"urls", newurlArray}}}})
	if err != nil {
		panic(err)
	}
}
