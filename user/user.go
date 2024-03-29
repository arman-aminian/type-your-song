package user

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(*model.User) error
	Remove(field, value string) error
	UpdateStrField(old *model.User, field string, value string) error
	UpdateStrFieldByEmail(old *model.User, field string, value string) error
	UpdateBoolField(old *model.User, field string, value bool) error
	UpdateBoolFieldByEmail(old *model.User, field string, value bool) error
	UpdateProfile(u *model.User) error

	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	GetById(primitive.ObjectID) (*model.User, error)

	AddFollowing(current primitive.ObjectID, u primitive.ObjectID) (model.User, error)
	RemoveFollowing(current primitive.ObjectID, u primitive.ObjectID) (model.User, error)

	Record(uid primitive.ObjectID, passed model.PassedSong, s *model.Song) (int, error)
	AddScore(uid primitive.ObjectID, score int) error
}
