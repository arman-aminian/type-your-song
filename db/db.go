package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

const (
	PATH = "mongodb://localhost:27017"
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(PATH)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		} else {
			err = client.Ping(context.TODO(), nil)
			if err != nil {
				clientInstanceError = err
			}
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func SetupUsersDb(mongoClient *mongo.Client) *mongo.Collection {
	usersDb := mongoClient.Database("type-your-song").Collection("users")
	createUniqueIndices(usersDb, "username")
	createUniqueIndices(usersDb, "email")
	return usersDb
}

func SetupSongsDb(mongoClient *mongo.Client) *mongo.Collection {
	songsDB := mongoClient.Database("type-your-song").Collection("songs")
	return songsDB
}

func SetupArtistsDb(mongoClient *mongo.Client) *mongo.Collection {
	artistsDB := mongoClient.Database("type-your-song").Collection("artists")
	return artistsDB
}

func SetupGenresDb(mongoClient *mongo.Client) *mongo.Collection {
	genresDB := mongoClient.Database("type-your-song").Collection("genres")
	return genresDB
}

func createUniqueIndices(db *mongo.Collection, field string) {
	_, err := db.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: field, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
