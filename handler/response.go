package handler

import (
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
)

type userResponse struct {
	User struct {
		Username       string `json:"username" bson:"_id"`
		Email          string `json:"email"`
		Name           string `json:"name"`
		Bio            string `json:"bio"`
		ProfilePicture string `json:"profile_picture"`
		Token          string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Name = u.Name
	r.User.Token = utils.GenerateJWT(u.ID.String())
	return r
}

type profileResponse struct {
	Profile struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"profile"`
}

func newProfileResponse(u *model.User) *profileResponse {
	r := new(profileResponse)
	r.Profile.Username = u.Username
	return r
}
