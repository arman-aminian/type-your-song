package song

import (
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(song *model.Song) error
	Remove(field, value string) error
	GetById(id primitive.ObjectID) (*model.Song, error)
}
