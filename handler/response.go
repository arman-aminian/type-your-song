package handler

import (
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
)

type userResponse struct {
	User struct {
		Username string `json:"username" bson:"_id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Token    string `json:"token"`
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
		Username string  `json:"username"`
		Image    *string `json:"image"`
	} `json:"profile"`
}

func newProfileResponse(u *model.User) *profileResponse {
	r := new(profileResponse)
	r.Profile.Username = u.Username
	return r
}
