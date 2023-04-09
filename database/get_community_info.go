package database

import (
	"context"
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCommunityInfo(gid uint32) (uint32, uint32, string, []uint32, []byte, uint64, uint64, uint32, uint32) {
	var result bson.M

	err := communitiesCollection.FindOne(context.TODO(), bson.D{{"gid", gid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}
	attribs := make([]uint32, len(result["attribs"].(bson.A)))
	for index, attrib := range result["attribs"].(bson.A) {
		if val, ok := attrib.(uint32); ok {
			attribs[index] = val
		}
	}

	application_buffer, _ := base64.StdEncoding.DecodeString(result["application_buffer"].(string))

	return uint32(result["host"].(int64)), uint32(result["type"].(int64)), result["password"].(string), attribs, application_buffer, uint64(result["start_date"].(int64)), uint64(result["end_date"].(int64)), uint32(result["sessions"].(int32)), uint32(len(result["participants"].(bson.A)))
}
