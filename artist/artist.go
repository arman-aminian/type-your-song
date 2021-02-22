package artist

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(*model.Artist) error
	Remove(string, string) error
	Find(primitive.ObjectID) (model.Artist, error)
}
