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
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
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
	co, err := utils.SaveToFiles(*cover, "files/songs/cover/", s.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on save song's cover")))
	}
	s.Cover = co

	gName := c.FormValue("genre")
	if gName == "" {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid genre")))
	}
	_, err = h.genreStore.GetByField("name", gName)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("genre not found"+err.Error())))
	}
	s.Genre = gName

	aID, err := primitive.ObjectIDFromHex(c.FormValue("artist"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("invalid artist")))
	}
	_, err = h.artistStore.Find(aID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("artist not found")))
	}
	s.Artist = aID

	s.Duration, err = strconv.Atoi(c.FormValue("duration"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid duration")))
	}
	s.PassedUsers = &[]primitive.ObjectID{}
	s.Url = c.FormValue("url")
	if s.Url == "" {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid url")))
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
	s.Score = utils.CalculateScore(s.WordsCount, s.MaxWPM, s.AvgWPM, utils.LevelToNum(utils.Hard))

	err = h.songStore.Create(&s)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on parse subtitle")))
	}

	err = h.genreStore.AddSong(s.ID, gName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on add song to genre")))
	}

	err = h.artistStore.AddSong(s.ID, aID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on add song to artist"+err.Error())))
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) DeleteSong(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
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

	songId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("invalid song id")))
	}
	s, err := h.songStore.GetById(songId)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("song not found")))
	}
	err = h.songStore.RemoveByID(songId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on remove song")))
	}

	err = h.genreStore.RemoveSong(s.ID, s.Genre)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on remove song from genre")))
	}

	err = h.artistStore.RemoveSong(s.ID, s.Artist)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on add song to artist"+err.Error())))
	}

	return c.JSON(http.StatusOK, s)
}

func (h *Handler) AddGenre(c echo.Context) error {
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

	var g model.Genre
	g.Name = c.FormValue("name")
	if g.Name == "" {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid name")))
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid cover")))
	}
	g.Cover, err = utils.SaveToFiles(*cover, "files/genres/cover/", g.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("could not save cover")))
	}
	g.Songs = &[]primitive.ObjectID{}
	err = h.genreStore.Create(&g)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("could not save genre")))
	}
	return c.JSON(http.StatusOK, g)
}

func (h *Handler) AddArtist(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
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

	var a model.Artist
	a.ID = primitive.NewObjectID()
	a.Name = c.FormValue("name")
	if a.Name == "" {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid name")))
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid cover")))
	}
	a.Cover, err = utils.SaveToFiles(*cover, "files/artists/cover/", a.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("could not save cover")))
	}
	a.Songs = &[]primitive.ObjectID{}
	err = h.artistStore.Create(&a)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("could not save artist")))
	}
	return c.JSON(http.StatusOK, a)
}

func (h *Handler) DeleteArtist(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
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

	aID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("invalid id")))
	}
	a, err := h.artistStore.Find(aID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("artist id not found")))
	}
	err = h.artistStore.RemoveByID(aID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("could not remove artist")))
	}
	for _, o := range *a.Songs {
		s, err := h.songStore.GetById(o)
		if err != nil {
			return c.JSON(http.StatusNotFound, utils.NewError(errors.New("song not found")))
		}
		err = h.songStore.RemoveByID(o)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("could not delete artist id : "+o.Hex())))
		}
		err = h.genreStore.RemoveSong(s.ID, s.Genre)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("error on remove song from genre")))
		}
	}
	return c.JSON(http.StatusOK, a)
}

func (h *Handler) GetSong(c echo.Context) error {
	jwtId := stringFieldFromToken(c, "id")
	sID, err := primitive.ObjectIDFromHex(c.Param("song"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("error on song id param")))
	}
	s, err := h.songStore.GetById(sID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	a, err := h.artistStore.Find(s.Artist)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	p := newFullSongResponse(s, a.Name)
	if jwtId != utils.Guest {
		id, err := primitive.ObjectIDFromHex(jwtId)
		if err == nil {
			cu, err := h.userStore.GetById(id)
			if err == nil {
				passed, err := utils.FindPassedSong(*cu.PassedSongs, s.ID)
				if err == nil {
					p.Song.PassedLevel = passed.PassedLevel
					p.Song.Speed = passed.Speed
					p.Song.Accuracy = passed.Accuracy
				}
			}
		}
	}
	return c.JSON(http.StatusOK, p)
}

func (h *Handler) GetSongs(c echo.Context) error {
	jwtId := stringFieldFromToken(c, "id")

	req := &songsIDRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	songs, err := h.songStore.GetSongs(req.Songs)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	loggedIn := false
	var response *fullSongsResponse
	if jwtId != utils.Guest {
		id, err := primitive.ObjectIDFromHex(jwtId)
		if err == nil {
			cu, err := h.userStore.GetById(id)
			if err == nil {
				loggedIn = true
				response = newFullSongsResponse(songs, h.artistStore, cu)
			}
		}
	}
	if !loggedIn {
		response = newFullSongsResponse(songs, h.artistStore, nil)
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetArtist(c echo.Context) error {
	aID, err := primitive.ObjectIDFromHex(c.Param("artist"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("error on artist id param")))
	}
	a, err := h.artistStore.Find(aID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newArtistResponse(&a))
}

func (h *Handler) GetGenre(c echo.Context) error {
	gN := c.Param("genre")
	if gN == "" {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("null genre")))
	}
	g, err := h.genreStore.GetByField("name", gN)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newGenreResponse(&g))
}

func (h *Handler) GetGenres(c echo.Context) error {
	g, err := h.genreStore.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newGenresResponse(&g))
}
