package models

type UserLogin struct {
	UsernameEmail string `json:"username_email" valid:"required"`
	Password      string `json:"password" valid:"required"`
}

type UserRegister struct {
	FirstName string `json:"first_name,omitempty" valid:"ascii,required"`
	LastName  string `json:"last_name,omitempty" valid:"ascii,required"`
	Email     string `json:"email,omitempty" valid:"email,required"`
	Username  string `json:"username,omitempty" valid:"ascii,required"`
	Password  string `json:"password,omitempty" valid:"ascii,required"`
	ConfirmPassword  string `json:"confirm_password,omitempty" valid:"ascii,required"`
	Role      string
}
type UserEdit struct {
	FirstName string `json:"first_name,omitempty" valid:"ascii,optional"`
	LastName  string `json:"last_name,omitempty"" valid:"ascii,optional"`
}

type UserEditEmail struct {
	NewEmail string `json:"email,omitempty" valid:"email,required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}