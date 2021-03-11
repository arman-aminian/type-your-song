package handler

import (
	"github.com/arman-aminian/type-your-song/router/middleware"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/labstack/echo"
)

func (h *Handler) Register(v1 *echo.Group) {
	guestUsers := v1.Group("/users")
	guestUsers.POST("/signup", h.SignUp)
	guestUsers.POST("/login", h.Login)
	guestUsers.GET("/login/google", h.GoogleLogin)
	guestUsers.GET("/callback", h.GoogleLoginCallback)
	guestUsers.POST("/reset", h.ResetPass)

	globalJwtMiddleware := middleware.JWTGlobal(utils.JWTSecret)
	globalUsers := v1.Group("", globalJwtMiddleware)
	globalUsers.GET("/user:username", h.GetProfile)
	globalUsers.GET("/song/:song", h.GetSong)
	globalUsers.GET("/songs", h.GetSongs)
	globalUsers.GET("/genre/:genre", h.GetGenre)
	globalUsers.GET("/genres", h.GetGenres)
	globalUsers.GET("/artist/:artist", h.GetArtist)

	confirmEmailJwtMiddleware := middleware.EmailConfirmJWT(utils.JWTSecret)
	confirmEmail := v1.Group("/confirm", confirmEmailJwtMiddleware)
	confirmEmail.GET("", h.ConfirmEmail)

	resetPassJwtMiddleware := middleware.ResetPassJWT(utils.JWTSecret)
	resetPass := v1.Group("/reset", resetPassJwtMiddleware)
	resetPass.GET("/confirm", h.ConfirmResetPass)

	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	users := v1.Group("/users", jwtMiddleware)
	users.POST("/follow/:username", h.Follow)
	users.DELETE("/unfollow/:username", h.UnFollow)

	songs := v1.Group("/song", jwtMiddleware)
	songs.POST("/add/song", h.AddSong)
	songs.DELETE("/delete/song/:id", h.DeleteSong)
	songs.POST("/add/genre", h.AddGenre)
	songs.POST("/add/artist", h.AddArtist)
	songs.DELETE("/delete/artist/:id", h.DeleteArtist)

	dummy := v1.Group("/dummy", jwtMiddleware)
	dummy.GET("", h.Dummy)

}
