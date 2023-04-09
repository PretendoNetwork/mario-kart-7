package database

import (
	"context"
	"time"

	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
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
	if globals.Config.DatabaseUseAuth {
		nexMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + globals.Config.DatabaseUsername + ":" + globals.Config.DatabasePassword + "@" + globals.Config.NEXDatabaseIP + ":" + globals.Config.NEXDatabasePort + "/"))
		accountMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + globals.Config.DatabaseUsername + ":" + globals.Config.DatabasePassword + "@" + globals.Config.AccountDatabaseIP + ":" + globals.Config.AccountDatabasePort + "/"))
	} else {
		nexMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + globals.Config.NEXDatabaseIP + ":" + globals.Config.NEXDatabasePort + "/"))
		accountMongoClient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://" + globals.Config.AccountDatabaseIP + ":" + globals.Config.AccountDatabasePort + "/"))
	}
	nexMongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = nexMongoClient.Connect(nexMongoContext)
	accountMongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = accountMongoClient.Connect(accountMongoContext)

	accountDatabase = accountMongoClient.Database(globals.Config.AccountDatabase)
	pnidCollection = accountDatabase.Collection(globals.Config.PNIDCollection)
	nexAccountsCollection = accountDatabase.Collection(globals.Config.NexAccountsCollection)

	splatoonDatabase = nexMongoClient.Database(globals.Config.SplatoonDatabase)
	regionsCollection = splatoonDatabase.Collection(globals.Config.RegionsCollection)
	usersCollection = splatoonDatabase.Collection(globals.Config.UsersCollection)
	sessionsCollection = splatoonDatabase.Collection(globals.Config.SessionsCollection)
	roomsCollection = splatoonDatabase.Collection(globals.Config.RoomsCollection)
	communitiesCollection = splatoonDatabase.Collection(globals.Config.CommunitiesCollection)

	sessionsCollection.DeleteMany(context.TODO(), bson.D{})
	roomsCollection.DeleteMany(context.TODO(), bson.D{})
}
