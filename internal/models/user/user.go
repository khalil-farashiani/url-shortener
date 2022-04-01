package user

import (
	"errors"
	"strings"

	"github.com/khalil-farashiani/url-shortener/internal/models/url"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint64    `json:"id" gorm:"primaryKey"`
	Username    string    `json:"user_name" gorm:"unique"`
	Email       *string   `json:"email"`
	Password    string    `json:"password"`
	Avatar      *string   `json:"avatar"`
	Phonenumber *string   `json:"phone_number"`
	Url         []url.Url `gorm:"foreignKey:UrlId"`
}

func (u *User) Validate() error {
	username := strings.TrimSpace(u.Username)
	password := strings.TrimSpace(u.Password)
	if username == "" {
		return errors.New("username should not be empty")
	}
	if len(password) < 8 {
		return errors.New("password character should be grater than or equal 8")
	}
	return nil
}
