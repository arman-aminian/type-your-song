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
	EmailConfirmJWTConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}
	EmailConfirmSkipper      func(c echo.Context) bool
	emailConfirmJwtExtractor func(echo.Context) (string, error)
)

var (
	EmailConfirmErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	EmailConfirmErrJWTInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt")
)

func EmailConfirmJWT(key interface{}) echo.MiddlewareFunc {
	c := EmailConfirmJWTConfig{}
	c.SigningKey = key
	return EmailConfirmJWTWithConfig(c)
}

func EmailConfirmJWTWithConfig(config EmailConfirmJWTConfig) echo.MiddlewareFunc {
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
				userID := claims["id"]
				c.Set("id", userID)
				userUN := claims["username"]
				c.Set("username", userUN)
				userEmail := claims["email"]
				fmt.Println("email :", userEmail)
				c.Set("email", userEmail)
				userPass := claims["password"]
				c.Set("password", userPass)

				return next(c)
			}
			return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
		}
	}
}
