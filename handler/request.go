package handler

import (
	"github.com/arman-aminian/type-your-song/model"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Registration request
type userRegisterRequest struct {
	User struct {
		Username string `json:"username" bson:"_id" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	u.ID = primitive.NewObjectID()
	return nil
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

type googleUserLoginRequest struct {
	//Content struct {
	Email   string `json:"email" validate:"required,email"`
	Picture string `json:"picture"`
	//} `json:"content"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type resetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func (r *resetPasswordRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

//************************ songs ************************
//type addSongRequest struct {
//	Url	string `json:"url"`
//
//}
