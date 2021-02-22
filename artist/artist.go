package artist

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(song *model.Artist) error
	Remove(field, value string) error
	Find(id primitive.ObjectID) (model.Artist, error)
}
