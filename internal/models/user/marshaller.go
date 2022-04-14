package user

import (
	"encoding/json"
)

type MarshallUser struct {
	Id        int64   `json:"id"`
	Username  string  `json:"username"`
	Email     *string `json:"email"`
	Avatar    *string `json:"avatar"`
	IsSpecial bool    `json:"is_special"`
}

func (u User) Marshall() MarshallUser {
	userJson, _ := json.Marshal(u)
	result := MarshallUser{}
	json.Unmarshal(userJson, &result)
	return result
}
