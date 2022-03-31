package url

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	UrlId    uint64 `json:"url_id" gorm:"constraint:OnDelete:SET NULL;"`
	Source   string `json:"source"`
	ShortUrl string `json:"short_url" gorm:"unique"`
}
