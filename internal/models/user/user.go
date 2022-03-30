package user

import (
	"errors"
	"strings"

	"github.com/khalil-farashiani/url-shortener/internal/models/url"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `json:"user_name"`
	Email    *string   `json:"email"`
	Password string    `json:"password"`
	Avatar   *string   `json:"avatar"`
	Url      []url.Url `gorm:"foreignKey:UrlId"`
}

func (u *User) Validate() error {
	username := strings.TrimSpace(u.Username)
	password := strings.TrimSpace(u.Password)
	if username == "" {
		return errors.New("username should not be empty")
	}
	if password == "" {
		return errors.New("password should not be empty")
	}
	return nil
}
