package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindCommunitiesWithParticipant(pid uint32) []uint32 {
	arr := []uint32{pid}
	gatheringIDs := []uint32{}

	cur, err := communitiesCollection.Find(context.TODO(), bson.M{"participants": bson.M{"$in": arr}}, options.Find())
	if err != nil {
		return nil
	}

	for cur.Next(context.TODO()) {
		fmt.Println("test")
		//Create a value into which the single document can be decoded
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return nil
		}

		gatheringIDs = append(gatheringIDs, (uint32)(result["gid"].(int64)))

	}

	return gatheringIDs
}
