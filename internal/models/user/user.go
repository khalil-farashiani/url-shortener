package user

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint64  `json:"id" gorm:"primaryKey"`
	Username    string  `json:"user_name" gorm:"unique"`
	Email       *string `json:"email"`
	Password    string  `json:"password"`
	Avatar      *string `json:"avatar"`
	Phonenumber *string `json:"phone_number"`
}

func (u *User) Validate(action string) error {
	switch action {
	case "save":
		username := strings.TrimSpace(u.Username)
		password := strings.TrimSpace(u.Password)
		if username == "" {
			return errors.New("username should not be empty")
		}
		if len(password) < 8 {
			return errors.New("password character should be grater than or equal 8")
		}
		return nil
	case "update":
		username := strings.TrimSpace(u.Username)
		password := strings.TrimSpace(u.Password)
		if username == "" {
			return errors.New("username should not be empty")
		}
		if len(password) < 8 {
			return errors.New("password character should be grater than or equal 8")
		}
		return nil
	default:
		return errors.New("action not in update or save")
	}
}
