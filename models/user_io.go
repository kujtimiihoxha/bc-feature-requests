package models

type UserLogin struct {
	UsernameEmail string `json:"username_email"`
	Password      string `json:"password"`
}

type UserRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ConfirmPassword  string `json:"confirm_password"`
	Role      string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}