package store

import (
	"context"
	"fmt"
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

func (gs *GenreStore) Remove(field, value string) error {
	_, err := gs.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (gs *GenreStore) Find(id primitive.ObjectID) (model.Genre, error) {
	var r model.Genre
	err := gs.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&r)
	return r, err
}

func (gs *GenreStore) Get(field, value string) (model.Genre, error) {
	var r model.Genre
	err := gs.db.FindOne(context.TODO(), bson.M{field: value}).Decode(&r)
	return r, err
}

func (gs *GenreStore) AddSong(sID primitive.ObjectID, to string) error {
	var err error
	g, err := gs.Get("name", to)
	if err != nil {
		return err
	}
	fmt.Println(g)
	fmt.Println(g.Songs)
	*g.Songs = append(*g.Songs, sID)
	fmt.Println(g)
	fmt.Println(g.Songs)
	_, err = gs.db.UpdateOne(context.TODO(), bson.M{"name": to}, bson.M{"$set": bson.M{"songs": g.Songs}})
	return err
}
