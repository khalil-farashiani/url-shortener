package url

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	UrlId    uint64 `json:"url_id"`
	Source   string `json:"source"`
	ShortUrl string `json:"short_url"`
}
