package user

import (
	"encoding/json"

	"github.com/khalil-farashiani/url-shortener/internal/models/url"
)

type MarshallUser struct {
	Id       int64     `json:"id"`
	Username string    `json:"user_name"`
	Email    *string   `json:"email"`
	Avatar   *string   `json:"avatar"`
	Url      []url.Url `json:"urls"`
}

func (user User) Marshall() MarshallUser {
	userJson, _ := json.Marshal(user)
	result := MarshallUser{}
	json.Unmarshal(userJson, &result)
	return result
}
