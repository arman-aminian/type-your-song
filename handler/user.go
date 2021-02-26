package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arman-aminian/type-your-song/email"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"strings"
)

func (h *Handler) SignUp(c echo.Context) error {
	var u model.User
	req := &userRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	t, err := h.userStore.GetByUsername(u.Username)
	if err == nil && t.HasPassword {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("duplicate username ")))
	}
	t, err = h.userStore.GetByEmail(u.Email)
	if err == nil && t.HasPassword {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("duplicate email ")))
	}

	emailJwt := utils.GenerateEmailConfirmJWT(u)
	to := []string{
		u.Email,
	}
	content := utils.BaseUrl + "/api/confirm?token=" + emailJwt
	err = email.SendEmail(to, content, "confirm your typeasong account")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusCreated, model.Message{Content: "an email sent to you\nconfirm your email address"})
}

func (h *Handler) ConfirmEmail(c echo.Context) error {
	var u model.User
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
		return err
	}
	u.ID = id
	u.Username = stringFieldFromToken(c, "username")
	u.Email = stringFieldFromToken(c, "email")
	u.Password = stringFieldFromToken(c, "password")
	u.HasPassword = true
	u.IsAdmin = false
	u.Followings = &[]primitive.ObjectID{}
	u.PassedSongs = &[]model.PassedSong{}
	u.Score = 0

	// todo error handling for duplicate click on confirm email
	if err := h.userStore.Create(&u); err != nil {
		if err = h.userStore.UpdateBoolFieldByEmail(&u, "has_password", true); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
		}
		if err = h.userStore.UpdateStrFieldByEmail(&u, "username", u.Username); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
		}
	}
	return c.JSON(http.StatusCreated, newUserResponse(&u))
}

func (h *Handler) Login(c echo.Context) error {
	req := &userLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u, err := h.userStore.GetByEmail(req.User.Email)
	if err != nil {
		return c.JSON(http.StatusForbidden, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	if !(u.CheckPassword(req.User.Password) && u.HasPassword) {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	return c.JSON(http.StatusOK, newUserResponse(u))
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/users/callback",
		ClientID:     "988996445182-6p95jh58g5kk1g0ecgn43gim0fnrvm40.apps.googleusercontent.com",
		ClientSecret: "vDdYMfEcajTBqKr6L0wqDESa",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

func (h *Handler) GoogleLogin(c echo.Context) error {
	token := utils.GenerateOauthToken()
	url := googleOauthConfig.AuthCodeURL(token)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleLoginCallback(c echo.Context) error {
	content, err := getUserInfo(c.FormValue("state"), c.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	var req googleUserLoginRequest
	err = json.Unmarshal(content, &req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	u, err := h.userStore.GetByEmail(req.Email)
	if err == nil && u != nil {
		return c.JSON(http.StatusOK, newUserResponse(u))
	}
	if err != nil {
		fmt.Println(err)
	}
	u.Email = req.Email
	u.Username = strings.Split(u.Email, "@")[0]
	u.ID = primitive.NewObjectID()
	u.HasPassword = false
	u.IsAdmin = false
	u.Followings = &[]primitive.ObjectID{}
	u.PassedSongs = &[]model.PassedSong{}
	u.Score = 0

	if err := h.userStore.Create(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newUserResponse(u))
}
func getUserInfo(state string, code string) ([]byte, error) {
	stateToken, err := jwt.Parse(state, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return utils.JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unexpected oauth state")
	}
	if !stateToken.Valid {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func (h *Handler) ResetPass(c echo.Context) error {
	e := c.QueryParam("email")
	fmt.Println("email :", e)
	u, err := h.userStore.GetByEmail(e)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NotFound())
	}
	//email := stringFieldFromToken(c, "email")

	emailJwt := utils.GenerateEmailConfirmJWT(*u)
	to := []string{
		e,
	}
	content := utils.BaseUrl + "/api/reset/confirm?token=" + emailJwt
	err = email.SendEmail(to, content, "reset password - typeasong")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("try again")))
	}
	return c.JSON(http.StatusCreated, model.Message{Content: "reset your password in the email we sent to you"})
}

func (h *Handler) ConfirmResetPass(c echo.Context) error {
	req := &resetPasswordRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	e := stringFieldFromToken(c, "email")
	u, err := h.userStore.GetByEmail(e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	hashedPass, err := u.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// todo error handling for duplicate click on confirm reset password
	if err := h.userStore.UpdateStrField(u, "password", hashedPass); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u.Password = hashedPass
	return c.JSON(http.StatusCreated, newUserResponse(u))
}

func (h *Handler) Follow(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.AccessForbidden())
	}
	cu, err := h.userStore.GetById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if cu == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if u.Username == cu.Username {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("can't follow yourself")))
	}
	if Contains(*cu.Followings, u.ID) {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("already follows the target")))
	}

	res, err := h.userStore.AddFollowing(cu.ID, u.ID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProfileResponse(&res))
}

func (h *Handler) UnFollow(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(stringFieldFromToken(c, "id"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.AccessForbidden())
	}
	cu, err := h.userStore.GetById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if cu == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if u.Username == cu.Username {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("can't unfollow yourself")))
	}
	if !Contains(*cu.Followings, u.ID) {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(errors.New("doesn't follow the target")))
	}

	res, err := h.userStore.RemoveFollowing(cu.ID, u.ID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	p := newProfileResponse(&res)
	return c.JSON(http.StatusOK, p)
}

func (h *Handler) GetProfile(c echo.Context) error {
	jwtId := stringFieldFromToken(c, "id")
	u, err := h.userStore.GetByUsername(c.Param("username"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	p := newProfileResponse(u)
	if jwtId != utils.Guest {
		id, err := primitive.ObjectIDFromHex(jwtId)
		if err == nil {
			cu, err := h.userStore.GetById(id)
			if err == nil {
				if Contains(*cu.Followings, u.ID) {
					p.Profile.IsFollowed = true
				}
			}
		}
	}
	return c.JSON(http.StatusOK, p)
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

func Contains(slice []primitive.ObjectID, val primitive.ObjectID) bool {
	if slice == nil {
		return true
	}
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
