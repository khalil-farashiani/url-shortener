package user

import (
	"encoding/json"
)

type MarshallUser struct {
	Id       int64   `json:"id"`
	Username string  `json:"user_name"`
	Email    *string `json:"email"`
	Avatar   *string `json:"avatar"`
}

func (user User) Marshall() MarshallUser {
	userJson, _ := json.Marshal(user)
	result := MarshallUser{}
	json.Unmarshal(userJson, &result)
	return result
}
