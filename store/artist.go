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

func (as *ArtistStore) RemoveByField(field, value string) error {
	_, err := as.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (as *ArtistStore) RemoveByID(id primitive.ObjectID) error {
	_, err := as.db.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (as *ArtistStore) Find(id primitive.ObjectID) (model.Artist, error) {
	var r model.Artist
	err := as.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&r)
	return r, err
}

func (as *ArtistStore) AddSong(sID primitive.ObjectID, to primitive.ObjectID) error {
	var err error
	a, err := as.Find(to)
	if err != nil {
		return err
	}
	*a.Songs = append(*a.Songs, sID)
	_, err = as.db.UpdateOne(context.TODO(), bson.M{"_id": to}, bson.M{"$set": bson.M{"songs": a.Songs}})
	return err
}

func (as *ArtistStore) RemoveSong(sID primitive.ObjectID, from primitive.ObjectID) error {
	var err error
	a, err := as.Find(from)
	if err != nil {
		return err
	}
	us := &[]primitive.ObjectID{}
	for _, o := range *a.Songs {
		if o != sID {
			*us = append(*us, o)
		}
	}
	_, err = as.db.UpdateOne(context.TODO(), bson.M{"_id": from}, bson.M{"$set": bson.M{"songs": us}})
	return err
}
