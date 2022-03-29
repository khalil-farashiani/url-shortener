package user

type User struct {
	Username string  `json:"user_name"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
	Avatar   *string `json:"avatar"`
}
