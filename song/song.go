package song

import "github.com/arman-aminian/type-your-song/model"

type Song interface {
	Create(song *model.Song) error
	Remove(field, value string) error
}
