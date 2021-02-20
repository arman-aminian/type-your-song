package model

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	HasPassword bool               `json:"has_password" bson:"has_password"`
	IsAdmin     bool               `json:"is_admin" bson:"is_admin"`
}

func NewUser() *User {
	var u User
	//u.Name = "Twitter User "
	//u.Tweets = &[]primitive.ObjectID{}
	//u.Followings = &[]Owner{}
	//u.Followers = &[]Owner{}
	//u.Notifications = &[]Event{}
	//u.Logs = &[]Event{}
	return &u
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
