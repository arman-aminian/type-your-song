package handler

import (
	"github.com/arman-aminian/type-your-song/artist"
	"github.com/arman-aminian/type-your-song/song"
	"github.com/arman-aminian/type-your-song/user"
)

type Handler struct {
	userStore   user.Store
	songStore   song.Store
	artistStore artist.Store
}

func NewHandler(us user.Store, ss song.Store, as artist.Store) *Handler {
	return &Handler{
		userStore:   us,
		songStore:   ss,
		artistStore: as,
	}
}
