package store

import (
	"context"
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	db *mongo.Collection
}

func NewUserStore(db *mongo.Collection) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) Create(u *model.User) error {
	_, err := us.db.InsertOne(context.TODO(), u)
	return err
}

func (us *UserStore) Remove(field, value string) error {
	_, err := us.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (us *UserStore) Update(old *model.User, field string, value string) error {
	var err error
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": old.ID}, bson.M{"$set": bson.M{field: value}})
	return err
}

func (us *UserStore) UpdateProfile(u *model.User) error {
	_, err := us.db.UpdateOne(context.TODO(),
		bson.M{"_id": u.Username},
		bson.M{"$set": bson.M{
			//"name": u.Name,
			//"bio":             u.Bio,
			//"profile_picture": u.ProfilePicture,
			//"header_picture":  u.HeaderPicture,
		},
		})
	return err
}

func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"username": username}).Decode(&u)
	return &u, err
}
