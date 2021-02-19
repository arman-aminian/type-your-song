package user

import (
	"github.com/arman-aminian/type-your-song/model"
)

type Store interface {
	Create(*model.User) error
	Remove(field, value string) error
	UpdateStrField(old *model.User, field string, value string) error
	UpdateStrFieldByEmail(old *model.User, field string, value string) error
	UpdateBoolField(old *model.User, field string, value bool) error
	UpdateBoolFieldByEmail(old *model.User, field string, value bool) error
	UpdateProfile(u *model.User) error

	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
}
