package handler

import (
	"errors"
	"fmt"
	ggosub "github.com/arman-aminian/gosub"
	"github.com/arman-aminian/gosub/parsers"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) AddSong(c echo.Context) error {
	println("user : " + stringFieldFromToken(c, "id"))
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnauthorized, utils.AccessForbidden())
	}
	u, err := h.userStore.GetById(id)
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
	lyrics, err := c.FormFile("lyric")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid srt file")))
	}
	f, err := os.Open(lyrics.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on parse subtitle")))
	}
	srt, err := ggosub.ParseByFile(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on parse subtitle")))
	}
	s.Lyrics = *srt
	s.WordsCount = parsers.TotalWordCount(srt)
	s.MaxWPM = parsers.CalculateMaxWpm(srt, 0, srt.Size)
	s.AvgWPM = parsers.CalculateMeanWpm(srt, 0, srt.Size)
	s.Score = calculateScore(s.WordsCount, s.MaxWPM, s.AvgWPM)

	//todo modify return value
	return c.JSON(http.StatusOK, s)
}
func calculateScore(size, max, avg int) int {
	t := (size / 1000) + (max / 300) + (avg / 200)
	return t / 3 * 100
}
