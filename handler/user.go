package handler

import (
	"fmt"
	"github.com/arman-aminian/type-your-song/email"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
	"net/http"
)

func (h *Handler) SignUp(c echo.Context) error {
	var u model.User
	req := &userRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	emailJwt := utils.GenerateEmailConfirmJWT(u)
	to := []string{
		u.Email,
	}
	err := email.SendEmail(to, emailJwt)
	if err != nil {
		panic(err)
	}

	if err := h.userStore.Create(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newUserResponse(&u))
}

func (h *Handler) ConfirmEmail(c echo.Context) error {
	id := stringFieldFromToken(c, "id")
	name := stringFieldFromToken(c, "name")
	username := stringFieldFromToken(c, "username")
	email := stringFieldFromToken(c, "email")
	pass := stringFieldFromToken(c, "password")

	return c.JSON(http.StatusCreated, email)
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
