package handler

import (
	"errors"
	"fmt"
	"github.com/arman-aminian/type-your-song/email"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *Handler) SignUp(c echo.Context) error {
	var u model.User
	req := &userRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	_, err := h.userStore.GetByUsername(u.Username)
	if err == nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("duplicate username ")))
	}
	_, err = h.userStore.GetByEmail(u.Email)
	if err == nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("duplicate email ")))
	}

	emailJwt := utils.GenerateEmailConfirmJWT(u)
	to := []string{
		u.Email,
	}
	content := utils.BaseUrl + "/api/confirm?query=" + emailJwt
	err = email.SendEmail(to, content)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusCreated, model.Message{Content: "an email sent to you\nconfirm your email address"})
}

func (h *Handler) ConfirmEmail(c echo.Context) error {
	var u model.User
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
		fmt.Println("here")
		return err
	}
	u.ID = id
	u.Name = stringFieldFromToken(c, "name")
	u.Username = stringFieldFromToken(c, "username")
	u.Email = stringFieldFromToken(c, "email")
	u.Password = stringFieldFromToken(c, "password")

	// todo error handling for duplicate click on confirm email
	if err := h.userStore.Create(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newUserResponse(&u))
}

func (h *Handler) Login(c echo.Context) error {
	username := c.Param("user")
	fmt.Println("id :", username)
	req := &userLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u, err := h.userStore.GetByEmail(req.User.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	if !u.CheckPassword(req.User.Password) {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	return c.JSON(http.StatusOK, newUserResponse(u))
}

func (h *Handler) Dummy(c echo.Context) error {
	return c.JSON(http.StatusCreated, "hello world")
}

func stringFieldFromToken(c echo.Context, field string) string {
	field, ok := c.Get(field).(string)
	if !ok {
		return ""
	}
	return field
}
