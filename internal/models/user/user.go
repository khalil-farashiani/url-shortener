package user

import (
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
