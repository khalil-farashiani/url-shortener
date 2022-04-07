package url

import (
	"encoding/json"
)

type MarshallUrl struct {
	Id       int64  `json:"id"`
	Source   string `json:"source"`
	ShortUrl string `json:"short_url" gorm:"unique"`
	UserID   uint64 `json:"user_id"`
}

func (u Url) Marshall() MarshallUrl {
	userJson, _ := json.Marshal(u)
	result := MarshallUrl{}
	json.Unmarshal(userJson, &result)
	return result
}

func (u Urls) Marshall() []interface{} {
	result := make([]interface{}, len(u))
	for index, url := range u {
		result[index] = url.Marshall()
	}
	return result
}
