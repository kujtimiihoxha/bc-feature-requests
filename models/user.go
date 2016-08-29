package models

/*
 User model.
*/

// Client model structure.
// FirstName: Client name.
// LastName: Client name.
// Email: Client name.
// Username: Client name.
// Password: Client name.
// Role: Additional data for clients
type User struct {
	BaseModel
	FirstName string `gorethink:"first_name" json:"first_name"`
	LastName  string `gorethink:"last_name" json:"last_name"`
	Email     string `gorethink:"email" json:"email"`
	Username  string `gorethink:"username" json:"username"`
	Password  string `gorethink:"password" json:"-"`
	Role      string `gorethink:"role" json:"role"`
}
