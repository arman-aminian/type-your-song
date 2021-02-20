package handler

import (
	"github.com/arman-aminian/type-your-song/router/middleware"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
)

func (h *Handler) Register(v1 *echo.Group) {
	guestUsers := v1.Group("/users")
	guestUsers.POST("", h.SignUp)
	guestUsers.POST("/login", h.Login)
	guestUsers.GET("/login/google", h.GoogleLogin)
	guestUsers.GET("/callback", h.GoogleLoginCallback)
	guestUsers.POST("/reset", h.ResetPass)

	confirmEmailJwtMiddleware := middleware.EmailConfirmJWT(utils.JWTSecret)
	confirmEmail := v1.Group("/confirm", confirmEmailJwtMiddleware)
	confirmEmail.GET("", h.ConfirmEmail)

	resetPassJwtMiddleware := middleware.ResetPassJWT(utils.JWTSecret)
	resetPass := v1.Group("/reset", resetPassJwtMiddleware)
	resetPass.GET("/confirm", h.ConfirmResetPass)

	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	songs := v1.Group("/song", jwtMiddleware)
	songs.POST("/add", h.AddSong)

	dummy := v1.Group("/dummy", jwtMiddleware)
	dummy.GET("", h.Dummy)

}
