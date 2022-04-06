package url

import (
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	Source   string    `json:"source"`
	ShortUrl string    `json:"short_url" gorm:"unique"`
	UserID   uint64    `json:"user_id"`
	User     user.User `jaon:"user"`
}
