package db

import (
	"github.com/astaxie/beego"
	r "github.com/dancannon/gorethink"
)

// Session holds the db session used from the models.
var session *r.Session

// Session used for tests.
var testSession *r.Mock

// Connect to the DB.
// If fail to connect there is no point to keep the server running.
// Connect must be called before the server is running (before beego.Run()).
func Connect() {
	s, err := r.Connect(r.ConnectOpts{
		Address:  beego.AppConfig.String("db::host"),
		Database: beego.AppConfig.String("db::database"),
	})
	if err != nil {
		beego.Error("Can not connect to the DB error : ", err)
		panic(err)
	}
	session = s
}
func Close() {
	session.Close()
}
func GetSession() interface{} {
	if beego.BConfig.RunMode == "test" {
		return testSession
	}
	return session
}
func SetTestSession(mock *r.Mock) {
	testSession = mock
}
