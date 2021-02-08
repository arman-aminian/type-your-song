package handler

import (
	"github.com/labstack/echo"
)

func (h *Handler) Register(v1 *echo.Group) {
	//jwtMiddleware := middleware.JWT(utils.JWTSecret)
	guestUsers := v1.Group("/users")
	guestUsers.POST("", h.SignUp)
	guestUsers.POST("/login", h.Login)

}
