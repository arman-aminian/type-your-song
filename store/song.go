package store

import (
	"context"
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type SongStore struct {
	db *mongo.Collection
}

func NewSongStore(db *mongo.Collection) *SongStore {
	return &SongStore{
		db: db,
	}
}

func (ss *SongStore) Create(s *model.Song) error {
	_, err := ss.db.InsertOne(context.TODO(), s)
	return err
}

func (ss *SongStore) RemoveByField(field, value string) error {
	_, err := ss.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (ss *SongStore) RemoveByID(id primitive.ObjectID) error {
	_, err := ss.db.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (ss *SongStore) GetById(id primitive.ObjectID) (*model.Song, error) {
	var s model.Song
	err := ss.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&s)
	return &s, err
}

func (ss *SongStore) GetSongs(ids []primitive.ObjectID) (*[]model.Song, error) {
	var result []model.Song
	query := bson.M{"_id": bson.M{"$in": ids}}
	res, err := ss.db.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = res.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
