package url

import (
	"github.com/khalil-farashiani/url-shortener/internal/models/user"
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	UrlId    uint64    `json:"url_id" gorm:"primaryKey;auto_increment;constraint:OnDelete:SET NULL;"`
	Source   string    `json:"source"`
	ShortUrl string    `json:"short_url" gorm:"unique"`
	UserID   uint64    `json:"user_id"`
	User     user.User `jaon:"user"`
}
