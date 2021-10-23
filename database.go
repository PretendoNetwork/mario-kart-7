package main

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContext context.Context
var accountDatabase *mongo.Database
var mk8Database *mongo.Database
var pnidCollection *mongo.Collection
var nexAccountsCollection *mongo.Collection
var regionsCollection *mongo.Collection
var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection
var roomsCollection *mongo.Collection
var tourneysCollection *mongo.Collection

func connectMongo() {
	if config.DatabaseUseAuth {
		mongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.DatabaseUsername + ":" + config.DatabasePassword + "@" + config.DatabaseIP + ":" + config.DatabasePort + "/"))
	} else {
		mongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.DatabaseIP + ":" + config.DatabasePort + "/"))
	}
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	accountDatabase = mongoClient.Database(config.AccountDatabase)
	pnidCollection = accountDatabase.Collection(config.PNIDCollection)
	nexAccountsCollection = accountDatabase.Collection(config.NexAccountsCollection)

	mk8Database = mongoClient.Database(config.MK8Database)
	regionsCollection = mk8Database.Collection(config.RegionsCollection)
	usersCollection = mk8Database.Collection(config.UsersCollection)
	sessionsCollection = mk8Database.Collection(config.SessionsCollection)
	roomsCollection = mk8Database.Collection(config.RoomsCollection)
	tourneysCollection = mk8Database.Collection(config.TournamentsCollection)

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

func addPlayerSession(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.InsertOne(context.TODO(), bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}})
	if err != nil {
		panic(err)
	}
}

func doesSessionExist(pid uint32) bool {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
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

func updatePlayerSessionAll(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}}}})
	if err != nil {
		panic(err)
	}
}

func updatePlayerSessionUrl(pid uint32, oldurl string, newurl string) {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
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

	_, err = sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"urls", newurlArray}}}})
	if err != nil {
		panic(err)
	}
}

func deletePlayerSession(pid uint32) {
	_, err := sessionsCollection.DeleteOne(context.TODO(), bson.D{{"pid", pid}})
	if err != nil {
		panic(err)
	}
}

func getPlayerUrls(pid uint32) []string {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	oldurlArray := result["urls"].(bson.A)
	newurlArray := make([]string, len(oldurlArray))
	for i := 0; i < len(oldurlArray); i++ {
		newurlArray[i] = oldurlArray[i].(string)
	}

	return newurlArray
}

func getPlayerSessionAddress(pid uint32) string {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	return result["ip"].(string) + ":" + result["port"].(string)
}
