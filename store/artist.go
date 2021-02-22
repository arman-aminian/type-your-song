package store

import (
	"context"
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArtistStore struct {
	db *mongo.Collection
}

func NewArtistStore(db *mongo.Collection) *ArtistStore {
	return &ArtistStore{
		db: db,
	}
}

func (as *ArtistStore) Create(s *model.Artist) error {
	_, err := as.db.InsertOne(context.TODO(), s)
	return err
}

func (as *ArtistStore) Remove(field, value string) error {
	_, err := as.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (as *ArtistStore) Find(id primitive.ObjectID) (model.Artist, error) {
	var r model.Artist
	err := as.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&r)
	return r, err
}
