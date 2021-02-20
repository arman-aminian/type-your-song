package handler

import (
	"errors"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

func (h *Handler) AddSong(c echo.Context) error {
	u, err := h.userStore.GetByUsername(stringFieldFromToken(c, "id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if !u.IsAdmin {
		return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("need admin permission")))
	}

	var s model.Song
	s.Name = c.FormValue("name")
	s.Cover = c.FormValue("cover")
	s.Genre = c.FormValue("genre")
	s.Artist, err = primitive.ObjectIDFromHex(c.FormValue("artist"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("artist not found")))
	}
	s.Duration, _ = strconv.Atoi(c.FormValue("duration"))
	s.PassedUsers = &[]primitive.ObjectID{}
	s.Url = c.FormValue("url")
	//ppf, err := c.FormFile("lyric")
	//todo modify return value
	return c.JSON(http.StatusOK, "added")
}
