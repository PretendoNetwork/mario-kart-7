package main

import (
	"context"
	"strconv"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var nexMongoClient *mongo.Client
var accountMongoClient *mongo.Client
var nexMongoContext context.Context
var accountMongoContext context.Context
var accountDatabase *mongo.Database
var splatoonDatabase *mongo.Database
var pnidCollection *mongo.Collection
var nexAccountsCollection *mongo.Collection
var regionsCollection *mongo.Collection
var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection
var roomsCollection *mongo.Collection
var communitiesCollection *mongo.Collection

func connectMongo() {
	if config.DatabaseUseAuth {
		nexMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.DatabaseUsername + ":" + config.DatabasePassword + "@" + config.NEXDatabaseIP + ":" + config.NEXDatabasePort + "/"))
		accountMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.DatabaseUsername + ":" + config.DatabasePassword + "@" + config.AccountDatabaseIP + ":" + config.AccountDatabasePort + "/"))
	} else {
		nexMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.NEXDatabaseIP + ":" + config.NEXDatabasePort + "/"))
		accountMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.AccountDatabaseIP + ":" + config.AccountDatabasePort + "/"))
	}
	nexMongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = nexMongoClient.Connect(nexMongoContext)
	accountMongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = accountMongoClient.Connect(accountMongoContext)

	accountDatabase = accountMongoClient.Database(config.AccountDatabase)
	pnidCollection = accountDatabase.Collection(config.PNIDCollection)
	nexAccountsCollection = accountDatabase.Collection(config.NexAccountsCollection)

	splatoonDatabase = nexMongoClient.Database(config.SplatoonDatabase)
	regionsCollection = splatoonDatabase.Collection(config.RegionsCollection)
	usersCollection = splatoonDatabase.Collection(config.UsersCollection)
	sessionsCollection = splatoonDatabase.Collection(config.SessionsCollection)
	roomsCollection = splatoonDatabase.Collection(config.RoomsCollection)
	communitiesCollection = splatoonDatabase.Collection(config.CommunitiesCollection)

	sessionsCollection.DeleteMany(context.TODO(), bson.D{})
	roomsCollection.DeleteMany(context.TODO(), bson.D{})
}

func getUsernameFromPID(pid uint32) string {
	var result bson.M

	err := pnidCollection.FindOne(context.TODO(), bson.D{{Key: "pid", Value: pid}}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ""
		}

		panic(err)
	}

	return result["username"].(string)
}

func isAllowedGameModeAndRegion(gamemode uint32, region uint32) bool {
	var result bson.M

	err := regionsCollection.FindOne(context.TODO(), bson.D{{Key: "id", Value: region}}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
	}

	if len(result["allowed_gamemodes"].(bson.A)) <= int(gamemode) {
		return false
	}
	if result["allowed_gamemodes"].(bson.A)[gamemode].(bool) {
		return true
	} else {
		return false
	}
}

func addNewUser(pid uint32) {
	_, err := usersCollection.InsertOne(context.TODO(), bson.D{{"pid", pid}, {"rating", bson.A{1000, 1000, 1000}}, {"username", getUsernameFromPID(pid)}, {"status", "allowed"}})
	if err != nil {
		panic(err)
	}
}

func updateUserRating(pid uint32, gamemode uint32, rating uint32) {
	_, err := usersCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"rating." + strconv.FormatUint(uint64(gamemode), 10), rating}}}})
	if err != nil {
		panic(err)
	}
}

func doesUserExist(pid uint32) bool {
	var result bson.M

	err := usersCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}

func addPlayerSession(rvcid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.InsertOne(context.TODO(), bson.D{{"rvcid", rvcid}, {"urls", urls}, {"ip", ip}, {"port", port}})
	if err != nil {
		panic(err)
	}
}

func doesSessionExist(rvcid uint32) bool {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"rvcid", rvcid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}

func updatePlayerSessionAll(rvcid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"rvcid", rvcid}}, bson.D{{"$set", bson.D{{"rvcid", rvcid}, {"urls", urls}, {"ip", ip}, {"port", port}}}})
	if err != nil {
		panic(err)
	}
}

func updatePlayerSessionUrl(rvcid uint32, oldurl string, newurl string) {
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

func deletePlayerSession(rvcid uint32) {
	_, err := sessionsCollection.DeleteOne(context.TODO(), bson.D{{"rvcid", rvcid}})
	if err != nil {
		panic(err)
	}
}

func getPlayerUrls(rvcid uint32) []string {
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

func getPlayerSessionAddress(rvcid uint32) string {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"rvcid", rvcid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(rvcid)
		//panic(err)
	}

	return result["ip"].(string) + ":" + result["port"].(string)
}
