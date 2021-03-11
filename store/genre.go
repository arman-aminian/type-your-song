package store

import (
	"context"
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GenreStore struct {
	db *mongo.Collection
}

func NewGenreStore(db *mongo.Collection) *GenreStore {
	return &GenreStore{
		db: db,
	}
}

func (gs *GenreStore) Create(s *model.Genre) error {
	_, err := gs.db.InsertOne(context.TODO(), s)
	return err
}

func (gs *GenreStore) RemoveByField(field, value string) error {
	_, err := gs.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (gs *GenreStore) RemoveByID(id primitive.ObjectID) error {
	_, err := gs.db.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (gs *GenreStore) Find(id primitive.ObjectID) (model.Genre, error) {
	var r model.Genre
	err := gs.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&r)
	return r, err
}

func (gs *GenreStore) GetByField(field, value string) (model.Genre, error) {
	var r model.Genre
	err := gs.db.FindOne(context.TODO(), bson.M{field: value}).Decode(&r)
	return r, err
}

func (gs *GenreStore) GetByID(id primitive.ObjectID) (model.Genre, error) {
	var r model.Genre
	err := gs.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&r)
	return r, err
}

func (gs *GenreStore) GetAll() ([]model.Genre, error) {
	var r []model.Genre
	cursor, err := gs.db.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &r); err != nil {
		return nil, err
	}
	return r, err
}

func (gs *GenreStore) AddSong(sID primitive.ObjectID, to string) error {
	var err error
	g, err := gs.GetByField("name", to)
	if err != nil {
		return err
	}
	*g.Songs = append(*g.Songs, sID)
	_, err = gs.db.UpdateOne(context.TODO(), bson.M{"name": to}, bson.M{"$set": bson.M{"songs": g.Songs}})
	return err
}

func (gs *GenreStore) RemoveSong(sID primitive.ObjectID, from string) error {
	var err error
	g, err := gs.GetByField("name", from)
	if err != nil {
		return err
	}
	us := &[]primitive.ObjectID{}
	for _, o := range *g.Songs {
		if o != sID {
			*us = append(*us, o)
		}
	}
	_, err = gs.db.UpdateOne(context.TODO(), bson.M{"name": from}, bson.M{"$set": bson.M{"songs": us}})
	return err
}
