package handler

import (
	"github.com/arman-aminian/type-your-song/user"
)

type Handler struct {
	userStore user.Store
	//articleStore article.Store
}

func NewHandler(us user.Store) *Handler {
	return &Handler{
		userStore: us,
		//articleStore: as,
	}
}
