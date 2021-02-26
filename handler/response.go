package handler

import (
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
