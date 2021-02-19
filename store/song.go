package store

import (
	"context"
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type SongStore struct {
	db *mongo.Collection
}

func NewSongStore(db *mongo.Collection) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (ss *SongStore) Create(s *model.Song) error {
	_, err := ss.db.InsertOne(context.TODO(), s)
	return err
}

func (ss *SongStore) Remove(field, value string) error {
	_, err := ss.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}
