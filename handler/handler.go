package handler

import (
	"github.com/arman-aminian/type-your-song/song"
	"github.com/arman-aminian/type-your-song/user"
)

type Handler struct {
	userStore user.Store
	songStore song.Store
}

func NewHandler(us user.Store, ss song.Store) *Handler {
	return &Handler{
		userStore: us,
		songStore: ss,
	}
}
