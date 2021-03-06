package artist

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(*model.Artist) error
	RemoveByField(string, string) error
	RemoveByID(id primitive.ObjectID) error
	Find(primitive.ObjectID) (model.Artist, error)
	AddSong(primitive.ObjectID, primitive.ObjectID) error
	RemoveSong(primitive.ObjectID, primitive.ObjectID) error
}
