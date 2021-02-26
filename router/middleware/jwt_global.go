package middleware

import (
	"fmt"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	JWTGlobalConfig struct {
		Skipper    SkipperGlobal
		SigningKey interface{}
	}
	SkipperGlobal      func(c echo.Context) bool
	jwtGlobalExtractor func(echo.Context) (string, error)
)

func JWTGlobal(key interface{}) echo.MiddlewareFunc {
	c := JWTGlobalConfig{}
	c.SigningKey = key
	return JWTGlobalWithConfig(c)
}

func JWTGlobalWithConfig(config JWTGlobalConfig) echo.MiddlewareFunc {
	extractor := jwtGlobalFromHeader("Authorization", "Token")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, err := extractor(c)
			if err != nil {
				if config.Skipper != nil {
					if config.Skipper(c) {
						return next(c)
					}
				}
				c.Set("id", utils.Guest)
				return next(c)
			}
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					fmt.Println("1")
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.SigningKey, nil
			})
			if err != nil {
				c.Set("id", utils.Guest)
				return next(c)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := claims["id"]
				c.Set("id", userID)
				fmt.Println("3")
				return next(c)
			}
			c.Set("id", utils.Guest)
			return next(c)
		}
	}
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtGlobalFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}
