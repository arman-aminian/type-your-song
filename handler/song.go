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
	"io/ioutil"
	"net/http"
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
	s.ID = primitive.NewObjectID()
	s.Name = c.FormValue("name")
	if s.Name == "" {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid name")))
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid cover")))
	}
	s.Cover, err = utils.SaveToFiles(*cover, "files/cover/", s.ID.Hex())

	s.Genre = c.FormValue("genre")
	// todo check genre is valid like artist!
	if s.Genre == "" {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid genre")))
	}
	s.Artist, err = primitive.ObjectIDFromHex(c.FormValue("artist"))
	if err != nil {
		//return c.JSON(http.StatusNotFound, utils.NewError(errors.New("artist not found")))
	}
	s.Duration, err = strconv.Atoi(c.FormValue("duration"))
	if err != nil {
		//return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid duration")))
	}
	s.PassedUsers = &[]primitive.ObjectID{}
	s.Url = c.FormValue("url")
	if s.Url == "" {
		//return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid url")))
	}
	lyrics, err := c.FormFile("lyric")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid lyric file")))
	}
	fh, err := lyrics.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on parse subtitle")))
	}
	f, err := ioutil.ReadAll(fh)
	srt, err := ggosub.ParseByFile(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on parse subtitle")))
	}
	s.Lyrics = *srt
	s.WordsCount = parsers.TotalWordCount(srt)
	s.MaxWPM = parsers.CalculateMaxWpm(srt, 0, srt.Size)
	s.AvgWPM = parsers.CalculateMeanWpm(srt, 0, srt.Size)
	s.Score = calculateScore(s.WordsCount, s.MaxWPM, s.AvgWPM)

	err = h.songStore.Create(&s)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on parse subtitle")))
	}

	return c.JSON(http.StatusOK, s)
}
func calculateScore(size, max, avg int) int {
	t := float64(size/1000) + float64(max/300) + float64(avg/200)
	return int(t / 3 * 100)
}
