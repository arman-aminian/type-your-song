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

func (gs *GenreStore) Create(s *model.Artist) error {
	_, err := gs.db.InsertOne(context.TODO(), s)
	return err
}

func (gs *GenreStore) Remove(field, value string) error {
	_, err := gs.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (gs *GenreStore) Find(id primitive.ObjectID) (model.Genre, error) {
	var r model.Genre
	err := gs.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&r)
	return r, err
}
