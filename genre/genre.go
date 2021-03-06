package genre

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(*model.Genre) error
	RemoveByField(string, string) error
	RemoveByID(id primitive.ObjectID) error
	Find(primitive.ObjectID) (model.Genre, error)
	GetByField(string, string) (model.Genre, error)
	GetByID(id primitive.ObjectID) (model.Genre, error)
	AddSong(primitive.ObjectID, string) error
	RemoveSong(primitive.ObjectID, string) error
}
