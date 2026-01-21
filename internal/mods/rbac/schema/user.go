package schema

import (
	"gin-admin/internal/config"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/errors"
	"gin-admin/pkg/util"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	UserStatusActivated = "activated"
	UserStatusFreezed   = "freezed"
)

type User struct {
	ID        string
	Username  string
	Name      string
	Password  string
	Phone     string
	Email     string
	Remark    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Roles     UserRoles
}

func (a *User) TableName() string {
	return config.C.FormatTableName("user")
}

type UserQueryParam struct {
	util.PaginationParam
	LikeUsername string
	LikeName     string
	Status       string
}

type UserQueryOptions struct {
	util.QueryOptions
}

type UserQueryResult struct {
	Data       Users
	PageResult *util.PaginationResult
}

type Users []*User

func (a Users) ToIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.ID)
	}
	return ids
}

type UserForm struct {
	Username string
	Name     string
	Password string
	Phone    string
	Email    string
	Remark   string
	Status   string
	Roles    UserRoles
}

func (a *UserForm) Validate() error {
	if a.Email != "" && validator.New().Var(a.Email, "email") != nil {
		return errors.BadRequest("", "Invalid email")
	}
	return nil
}

func (a *UserForm) FillTo(user *User) error {
	user.Username = a.Username
	user.Name = a.Name
	user.Phone = a.Phone
	user.Email = a.Email
	user.Remark = a.Remark
	user.Status = a.Status

	if pass := a.Password; pass != "" {
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return errors.BadRequest("", "GeneratePassword %s error", err.Error())
		}
		user.Password = hashPass
	}
	return nil
}
