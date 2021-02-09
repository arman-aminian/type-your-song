package utils

import (
	"github.com/arman-aminian/type-your-song/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTSecret = []byte("!!SECRET!!")

func GenerateJWT(id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}

func GenerateEmailConfirmJWT(user model.User) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID.Hex()
	claims["name"] = user.Name
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["password"] = user.Password
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}
