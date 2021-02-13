package user

import (
	"github.com/arman-aminian/type-your-song/model"
)

type Store interface {
	Create(*model.User) error
	Remove(field, value string) error
	Update(old *model.User, field string, value string) error
	UpdateProfile(u *model.User) error

	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	//AddFollower(u *model.User, follower *model.User) error
	//RemoveFollower(u *model.User, follower *model.User) error
	//IsFollower(username, followerUsername string) (bool, error)
	//
	//AddTweet(u *model.User, t *model.Tweet) error
	//RemoveTweet(u *model.User, id *string) error
	//
	//AddLog(u *model.User, e *model.Event) error
	//AddNotification(u *model.User, e *model.Event) error
	//
	//GetUserListFromUsernameList(usernames []string) (*[]model.User, error)
	//GetTweetIdListFromUsernameList(usernames []string) (*[]primitive.ObjectID, error)
	//
	//GetUsernameSearchResult(username string) (*[]model.Owner, error)
}
