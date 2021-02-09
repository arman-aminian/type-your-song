package middleware

import (
	"errors"
	"fmt"
	"github.com/arman-aminian/type-your-song/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	ResetPassJWTConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}
	ResetPassSkipper      func(c echo.Context) bool
	ResetPassJwtExtractor func(echo.Context) (string, error)
)

var (
	ResetPassErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	ResetPassErrJWTInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt")
)

func ResetPassJWT(key interface{}) echo.MiddlewareFunc {
	c := ResetPassJWTConfig{}
	c.SigningKey = key
	return ResetPassJWTWithConfig(c)
}

func ResetPassJWTWithConfig(config ResetPassJWTConfig) echo.MiddlewareFunc {
	//extractor := jwtFromHeader("Authorization", "Token")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.QueryParam("token")
			println("auth :", auth)
			if auth == "" {
				if config.Skipper != nil {
					if config.Skipper(c) {
						return next(c)
					}
				}
				return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("not found! ")))
			}
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.SigningKey, nil
			})
			if err != nil {
				return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userEmail := claims["email"]
				c.Set("email", userEmail)
				return next(c)
			}
			return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
		}
	}
}
