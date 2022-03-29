package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `json:"user_name"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
	Avatar   *string `json:"avatar"`
}
