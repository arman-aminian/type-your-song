package handler

import (
	"github.com/arman-aminian/gosub/parsers"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userResponse struct {
	User struct {
		Username string `json:"username" bson:"_id"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Token = utils.GenerateJWT(u.ID.Hex())
	return r
}

type profileResponse struct {
	Profile struct {
		Username    string                `json:"username"`
		Email       string                `json:"email"`
		Image       string                `json:"image"`
		PassedSongs *[]model.PassedSong   `json:"passed_songs"`
		Followings  *[]primitive.ObjectID `json:"followings"`
		Score       int                   `json:"score"`
		IsFollowed  bool                  `json:"is_followed"`
	} `json:"profile"`
}

func newProfileResponse(u *model.User) *profileResponse {
	r := new(profileResponse)
	r.Profile.Username = u.Username
	r.Profile.Email = u.Email
	r.Profile.Image = u.Image
	r.Profile.PassedSongs = u.PassedSongs
	r.Profile.Followings = u.Followings
	r.Profile.Score = u.Score
	r.Profile.IsFollowed = false
	return r
}

type fullSongResponse struct {
	Song struct {
		ID          primitive.ObjectID    `json:"id"`
		Name        string                `json:"name"`
		Cover       string                `json:"cover"`
		Genre       string                `json:"genre"`
		ArtistID    primitive.ObjectID    `json:"artist_id"`
		ArtistName  string                `json:"artist_name"`
		Duration    int                   `json:"duration"`
		PassedUsers *[]primitive.ObjectID `json:"passed_users"`
		Url         string                `json:"url"`
		Lyrics      parsers.Srt           `json:"lyrics"`

		MaxWPM     int `json:"max_wpm"`
		AvgWPM     int `json:"avg_wpm"`
		WordsCount int `json:"words_count"`
		Score      int `json:"score"`

		PassedLevel string `json:"passed_level"`
		Speed       int    `json:"speed"`
		Accuracy    int    `json:"accuracy"`
	}
}

func newFullSongResponse(s *model.Song, artist string) *fullSongResponse {
	r := new(fullSongResponse)
	r.Song.ID = s.ID
	r.Song.Name = s.Name
	r.Song.Cover = s.Cover
	r.Song.Genre = s.Genre
	r.Song.ArtistID = s.Artist
	r.Song.ArtistName = artist
	r.Song.Duration = s.Duration
	r.Song.PassedUsers = s.PassedUsers
	r.Song.Url = s.Url
	r.Song.Lyrics = s.Lyrics
	r.Song.MaxWPM = s.MaxWPM
	r.Song.AvgWPM = s.AvgWPM
	r.Song.WordsCount = s.WordsCount
	r.Song.Score = s.Score
	return r
}

type genreResponse struct {
	Genre struct {
		Name  string                `json:"name"`
		Cover string                `json:"cover"`
		Songs *[]primitive.ObjectID `json:"songs"`
	}
}

func newGenreResponse(g *model.Genre) *genreResponse {
	r := new(genreResponse)
	r.Genre.Songs = g.Songs
	r.Genre.Name = g.Name
	r.Genre.Cover = g.Cover
	return r
}

type genresResonse struct {
	Genres []model.Genre `json:"genres"`
}

func newGenresResponse(g *[]model.Genre) *genresResonse {
	r := new(genresResonse)
	r.Genres = *g
	return r
}

type artistResponse struct {
	Artist struct {
		ID    primitive.ObjectID    `json:"_id"`
		Name  string                `json:"name"`
		Cover string                `json:"cover"`
		Songs *[]primitive.ObjectID `json:"songs"`
	}
}

func newArtistResponse(a *model.Artist) *artistResponse {
	r := new(artistResponse)
	r.Artist.Cover = a.Cover
	r.Artist.Name = a.Name
	r.Artist.Songs = a.Songs
	r.Artist.ID = a.ID
	return r
}
