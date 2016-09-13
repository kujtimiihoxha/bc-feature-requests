package models

import (
	r "github.com/dancannon/gorethink"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
	"github.com/astaxie/beego"
	"strings"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
)
/*
 User model.
*/

// User model structure.
// FirstName: User first name.
// LastName: User name.
// Email: User email.
// Username: User username.
// Password: User password.
// Role: User role
type User struct {
	BaseModel
	ClientName string `gorethink:"client_name,omitempty" json:"client_name"`
	ClientID string `gorethink:"client_id,omitempty" json:"client_id"`
	FirstName string `gorethink:"first_name,omitempty" json:"first_name"`
	LastName  string `gorethink:"last_name,omitempty" json:"last_name"`
	Email     string `gorethink:"email,omitempty" json:"email"`
	Username  string `gorethink:"username,omitempty" json:"username"`
	Password  string `gorethink:"password,omitempty" json:"-"`
	Role      int    `gorethink:"role,omitempty" json:"role"`
	Verified      bool    `gorethink:"verified,omitempty" json:"-"`
}

const user_table = "users"

// Claims contains the claims that will be encoded in the token.
type Claims struct {
	jwt.StandardClaims
	Username  string  `json:"username,omitempty"`
	Email  string  `json:"email,omitempty"`
	FirstName string  `json:"firstname,omitempty"`
	ClientName string  `json:"clientname,omitempty"`
	ClientID string  `json:"clientId,omitempty"`
	LastName  string  `json:"lastname,omitempty"`
	Role      int  `json:"role,omitempty"`
	ID      string    `json:"id"`
}

const (
	ADMIN_ROLE = 1
	EMPLOY = 2
	CLIENT = 3
)

// Create new user from UserRegister data.
// user: data.
// t: time to set CreatedAt.
// Returns:
//	- The User created.
func NewUser(user *UserRegister, t time.Time) (*User, *CodeInfo) {
	u := &User{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: strings.ToLower(user.Email),
		Username: strings.ToLower(user.Username),
		Role: EMPLOY,
		Verified:false,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}

	res := u.validateRegistration(user)
	if res.Code != 0 {
		return nil, res
	}
	res = u.generateUserPassword(user)
	if res.Code != 0 {
		return nil, res
	}
	return u,OkInfo("")
}

// Create new user from UserRegister data.
// user: data.
// t: time to set CreatedAt.
// Returns:
//	- The User created.
func NewClientUser(user *ClientRegister, t time.Time) (*User,*Client, *CodeInfo) {
	u := &User{
		ClientName: user.ClientName,
		Email: strings.ToLower(user.Email),
		Username: strings.ToLower(user.Username),
		Role: CLIENT,
		Verified:false,
		BaseModel: BaseModel{
			CreatedAt: &t,
		},
	}
	clientCreate := &ClientCreate{
		Name:user.ClientName,
		Description:user.Description,
	}
	client := NewClient(clientCreate,time.Now().UTC())
	usr:= &UserRegister{
		Email:user.Email,
		Password:user.Password,
		ConfirmPassword:user.ConfirmPassword,
		Username:user.Username,
	}
	res := u.validateRegistration(usr)
	if res.Code != 0 {
		return nil,nil, res
	}
	res = u.generateUserPassword(usr)
	if res.Code != 0 {
		return  nil,nil, res
	}
	return u,client,OkInfo("")
}
// Get all users from the DB.
// Returns:
// 	- Array of users (or empty if there are no users in the DB).
// 	- CodeInfo with the error information.
func GetAllUsers() ([]User, *CodeInfo) {
	users := []User{}
	result := getAll(user_table, &users)
	return users, result
}
// Get all employs from the DB.
// Returns:
// 	- Array of users (or empty if there are no users in the DB).
// 	- CodeInfo with the error information.
func GetAllEmploys() ([]User, *CodeInfo) {
	users := []User{}
	res,err := r.Table(user_table).Filter(r.Row.Field("role").Ne(1)).Run(db.GetSession().(r.QueryExecutor))
	if err != nil {
		return users,  ErrorInfo(ErrSystem, err.Error())
	}
	err = res.All(&users)
	if err != nil {
		return users, ErrorInfo(ErrSystem, err.Error())
	}
	return users,  OkInfo("")
}

// Get user by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the data of the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (u *User) GetById(id string) *CodeInfo {
	return u.getById(user_table, id, u)
}

// Get user by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the data of the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (u *User) GetNotifications(id string) ([]Notification,*CodeInfo){
	notifications:=[]Notification{}
	notificationsRes,err := r.Table(notifications_table).Filter(r.And(r.Row.Field("user_id").Eq(id),r.Row.Field("viewed").Eq(false))).Run(u.Session())
	if err != nil {
		return notifications, ErrorInfo(ErrSystem, err.Error())
	}
	notificationsRes.All(&notifications)
	return notifications, OkInfo("")
}


func (u *User) SetNotificationsViewed(id string) *CodeInfo{
	wr,err := r.Table(notifications_table).Filter(r.And(r.Row.Field("user_id").Eq(id),r.Row.Field("viewed").Eq(false))).Update(Notification{
		Viewed:true,
	}).RunWrite(u.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	return OkInfo("")
}

// Insert new user.
// User should have data before calling this method.
// Error :
// 	- Returns CodeInfo with the error informatiogn.
// Success :
//     - Sets the ID of the model calling the method on Success
//     - Returns CodeInfo with Code = 0 (No error)

func (u *User) Insert() *CodeInfo {
	id, result := u.insert(user_table, u)
	u.ID = id
	return result
}
// Insert new client user.
// User should have data before calling this method.
// Error :
// 	- Returns CodeInfo with the error informatiogn.
// Success :
//     - Sets the ID of the model calling the method on Success
//     - Returns CodeInfo with Code = 0 (No error)

func (u *User) InsertClient(client *Client) *CodeInfo {
	idClient, result := client.insert(client_table, client)
	if result.Code != 0 {
		return  result
	}
	u.ClientID = idClient;
	id, result := u.insert(user_table, u)
	u.ID = id
	return result
}
func (u *User) Verify(id string, t time.Time) *CodeInfo {
	result:=u.GetById(id)
	if result.Code != 0 {
		return result
	}
	if u.Verified {
		return ErrorInfo(ErrUserAlreadyVerified,"User is already verified")
	}
	uv := User{}
	uv.UpdatedAt = &t
	uv.Verified = true
	return u.update(user_table,id,uv)
}

// Update user by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the updated data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (u *User) Update(username string, data *UserEdit, t time.Time) *CodeInfo {
	u.setUserFromEdit(data)
	u.UpdatedAt = &t
	wr, err := r.Table(user_table).Filter(r.Row.Field("username").Eq(username)).Update(u).RunWrite(u.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if wr.Errors > 0 {
		return ErrorInfo(ErrDatabase, wr.FirstError)
	}
	rs,err:= r.Table(user_table).Filter(r.Row.Field("username").Eq(username)).Run(u.Session());
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if rs.IsNil() {
		return  ErrorInfo(ErrNotFound, "Not Found")
	}
	rs.One(&u)
	return OkInfo("Data updated succesfully")
}

// Delete user by id.
// Error :
// 	- Returns CodeInfo with the error information.
// Success :
//     - Fills the deleted data to the model calling the method.
//     - Returns CodeInfo with Code = 0 (No error)
func (c *User) Delete(id string) *CodeInfo {
	return c.delete(user_table, id, c)
}

// Set User data from UserEdit model.
func (u *User) setUserFromEdit(data *UserEdit) {
	u.FirstName = data.FirstName
	u.LastName = data.LastName
}
// Set User data from UserEdit model.
func (u *User) generateUserPassword(user *UserRegister) *CodeInfo {
	res := u.validatePasswords(user);
	if res.Code != 0 {
		return res
	}
	d,err :=bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil {
		return ErrorInfo(ErrSystem,err.Error())
	}
	user.Password =""
	u.Password = string(d)
	return OkInfo("")
}
// Set User data from UserEdit model.
func (u *User) validatePasswords(user * UserRegister) *CodeInfo {
	if user.Password != user.ConfirmPassword {
		return ErrorInfo(ErrPasswordMismatch,"Passwords do not match")
	}
	return OkInfo("")
}
func (u *User) validateRegistration(user *UserRegister) *CodeInfo {
	usersWithEmail, err := r.Table(user_table).Filter(r.Row.Field("email").Eq(user.Email)).Run(u.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if !usersWithEmail.IsNil() {
		beego.Debug("Email already exists")
		return ErrorInfo(ErrEmailExists, "Email already exists")
	}
	usersWithUsername, err := r.Table(user_table).Filter(r.Row.Field("username").Eq(user.Username)).Run(u.Session())
	if err != nil {
		return ErrorInfo(ErrSystem, err.Error())
	}
	if !usersWithUsername.IsNil() {
		beego.Debug("Username already exists")
		return ErrorInfo(ErrUsernameExists, "Username already exists")
	}
	return OkInfo("")
}

// User login.
// Error :
// 	- Returns CodeInfo with the error informatiogn.
// Success :
//     - Returns TokenResponse and  CodeInfo with Code = 0 (No error)
func (u *User) Login(userLogin UserLogin) (*TokenResponse, *CodeInfo) {
	user := User{}
	userResponse, err := r.Table(user_table).Filter(
		r.Or(r.Row.Field("username").Eq(strings.ToLower(userLogin.UsernameEmail)),
			r.Row.Field("email").Eq(strings.ToLower(userLogin.UsernameEmail)))).Run(u.Session())
	if err != nil {
		return nil, ErrorInfo(ErrSystem, err.Error())
	}
	if userResponse.IsNil() {
		return nil, ErrorInfo(ErrNotFound, "User Not Found")
	}
	userResponse.One(&user)
	if !user.Verified {
		return nil, ErrorInfo(ErrUserNotVerified, "User is not verified")
	}
	err = bcrypt.CompareHashAndPassword(([]byte)(user.Password), ([]byte)(userLogin.Password))
	if err != nil {
		return nil, ErrorInfo(ErrUnAuthorized, "Password does not match")
	}
	v, _ := beego.AppConfig.Int("jwt::hours");
	claims := &Claims{
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour * time.Duration(v)).Unix(),
			Issuer:    "bc",
		},
		Username: user.Username,
		Email: user.Email,
		FirstName: user.FirstName,
		ClientName:user.ClientName,
		ClientID:user.ClientID,
		LastName: user.LastName,
		Role: user.Role,
		ID: user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(beego.AppConfig.String("jwt::key")))
	if err != nil {
		return nil, ErrorInfo(ErrSystem, err.Error())
	}
	tr := &TokenResponse{
		Token: ss,
	}
	return tr, OkInfo("Authenticated")
}
