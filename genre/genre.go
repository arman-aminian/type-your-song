package genre

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(*model.Genre) error
	Remove(string, string) error
	Find(primitive.ObjectID) (model.Genre, error)
	Get(string, string) (model.Genre, error)
}
