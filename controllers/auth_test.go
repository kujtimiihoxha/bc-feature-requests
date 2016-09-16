package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	r "github.com/dancannon/gorethink"
	"github.com/dgrijalva/jwt-go"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	_ "github.com/kujtimiihoxha/bc-feature-requests/tests"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var user_table = "users"

// Init client routes  for testing.
func init() {
	// Base route
	beego.Router("/", &MainController{})
	// Api V1 routes
	ns := beego.NewNamespace("/api/v1",
		// Authentication endpoint
		beego.NSNamespace("/auth",
			// Log in
			beego.NSRouter("/login", &AuthController{}, "post:Login"),
		),
	)
	// Add Api v1 namespace to beego.
	beego.AddNamespace(ns)
}

/*-------------------------------
	Test Login
--------------------------------*/

// Test the response when there is any database error.
// Should Return :
// Status : 200
// TokenResponse : With valid token
func TestUserLoginSuccessWithUsername(t *testing.T) {
	userLogin := models.UserLogin{
		UsernameEmail: "employ",
		Password:      "employ@123",
	}
	b, _ := json.Marshal(userLogin)
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	mock := r.NewMock()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: "585aca07-6d3c-43ba-97d4-8fb4cb27e024",
		},
		Username:  "employ",
		Email:     "employ@gmail.com",
		Role:      2,
		FirstName: "Employ",
		Verified:  true,
		LastName:  "Name",
		Password:  "$2a$10$Rtt0sfArkW1gCLeiW5AUbu6VgRNtzzYRKPmD5xmK/JhAyw4VA8Ipq",
	}
	db.SetTestSession(mock)
	mock.On(r.Table(user_table).Filter(
		r.Or(r.Row.Field("username").Eq(strings.ToLower(userLogin.UsernameEmail)),
			r.Row.Field("email").Eq(strings.ToLower(userLogin.UsernameEmail))))).Once().Return(user, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if login succesfull\n", t, func() {
		response := models.TokenResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Response tocken should have the same encoded role as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(int(token.Claims.(jwt.MapClaims)["role"].(float64)), ShouldEqual, user.Role)
		})
		Convey("Response tocken should have the same encoded username as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["username"], ShouldEqual, user.Username)
		})
		Convey("Response tocken should have the same encoded firstname as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["firstname"], ShouldEqual, user.FirstName)
		})
		Convey("Response tocken should have the same encoded lastname as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["lastname"], ShouldEqual, user.LastName)
		})
		Convey("Response tocken should have the same encoded id as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["id"], ShouldEqual, user.ID)
		})
	})
}
func TestUserLoginSuccessWithEmail(t *testing.T) {
	userLogin := models.UserLogin{
		UsernameEmail: "employ@gmail.com",
		Password:      "employ@123",
	}
	b, _ := json.Marshal(userLogin)
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	mock := r.NewMock()
	user := models.User{
		BaseModel: models.BaseModel{
			ID: "585aca07-6d3c-43ba-97d4-8fb4cb27e024",
		},
		Username:  "employ",
		Email:     "employ@gmail.com",
		Role:      2,
		FirstName: "Employ",
		LastName:  "Name",
		Verified:  true,
		Password:  "$2a$10$Rtt0sfArkW1gCLeiW5AUbu6VgRNtzzYRKPmD5xmK/JhAyw4VA8Ipq",
	}
	db.SetTestSession(mock)
	mock.On(r.Table(user_table).Filter(
		r.Or(r.Row.Field("username").Eq(strings.ToLower(userLogin.UsernameEmail)),
			r.Row.Field("email").Eq(strings.ToLower(userLogin.UsernameEmail))))).Once().Return(user, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if login succesfull\n", t, func() {
		response := models.TokenResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Response tocken should have the same encoded role as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(int(token.Claims.(jwt.MapClaims)["role"].(float64)), ShouldEqual, user.Role)
		})
		Convey("Response tocken should have the same encoded username as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["username"], ShouldEqual, user.Username)
		})
		Convey("Response tocken should have the same encoded firstname as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["firstname"], ShouldEqual, user.FirstName)
		})
		Convey("Response tocken should have the same encoded lastname as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["lastname"], ShouldEqual, user.LastName)
		})
		Convey("Response tocken should have the same encoded id as the user in the DB", func() {
			token, _ := jwt.Parse(response.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(beego.AppConfig.String("jwt::key")), nil
			})
			So(token.Claims.(jwt.MapClaims)["id"], ShouldEqual, user.ID)
		})
	})
}

// Tests if the update method fails
//  if no data is sent.
// Should Return :
// Status : 400
// Code : 10001
func TestUserLoginFailWhenNoData(t *testing.T) {
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	FailNoData(w, t)
}

// Tests if the update method fails
//  if the data send is not valid.
// Should Return :
// Status : 400
// Code : 10014
func TestLoginFailValidation(t *testing.T) {
	userLogin := models.UserLogin{
		UsernameEmail: "employ@gmail.com",
	}
	b, _ := json.Marshal(userLogin)
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	ValidationFail(w, t, "Password: non zero value required;")
}

// Test the response when there is any system error.
// Should Return :
// Status : 500
// Code : 10011
func TestLoginSystemError(t *testing.T) {
	userLogin := models.UserLogin{
		UsernameEmail: "employ@gmail.com",
		Password:      "employ@123",
	}
	b, _ := json.Marshal(userLogin)
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(user_table).Filter(
		r.Or(r.Row.Field("username").Eq(strings.ToLower(userLogin.UsernameEmail)),
			r.Row.Field("email").Eq(strings.ToLower(userLogin.UsernameEmail))))).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	SystemError(w, t)
}

// Test the response when the user could not be found.
// Should Return :
// Status : 404
// Code : 404
func TestUserNotFound(t *testing.T) {
	userLogin := models.UserLogin{
		UsernameEmail: "employ@gmail.com",
		Password:      "employ@123",
	}
	b, _ := json.Marshal(userLogin)
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(user_table).Filter(
		r.Or(r.Row.Field("username").Eq(strings.ToLower(userLogin.UsernameEmail)),
			r.Row.Field("email").Eq(strings.ToLower(userLogin.UsernameEmail))))).Once().Return(nil, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	NotFound(w, t, "User Not Found")

}

// Test the response when the user password is not matching
// Should Return :
// Status : 404
// Code : 404
func TestPasswordIncorrect(t *testing.T) {
	userLogin := models.UserLogin{
		UsernameEmail: "employ@gmail.com",
		Password:      "employ@1234",
	}
	b, _ := json.Marshal(userLogin)
	rs, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	user := models.User{
		BaseModel: models.BaseModel{
			ID: "585aca07-6d3c-43ba-97d4-8fb4cb27e024",
		},
		Username:  "employ",
		Email:     "employ@gmail.com",
		Role:      2,
		FirstName: "Employ",
		LastName:  "Name",
		Verified:  true,
		Password:  "$2a$10$Rtt0sfArkW1gCLeiW5AUbu6VgRNtzzYRKPmD5xmK/JhAyw4VA8Ipq",
	}
	mock.On(r.Table(user_table).Filter(
		r.Or(r.Row.Field("username").Eq(strings.ToLower(userLogin.UsernameEmail)),
			r.Row.Field("email").Eq(strings.ToLower(userLogin.UsernameEmail))))).Once().Return(user, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if not found is returned \n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Code should be 10013", func() {
			So(response.Code, ShouldEqual, 10013)
		})
		Convey("Message should specify the system erorr message", func() {
			So(response.Message, ShouldEqual, "Password does not match")
		})
	})

}

func TestMostBeAuthorizedAllowsOPTIONS(t *testing.T) {
	r := &http.Request{}
	c := context.NewContext()
	c.Request = r
	r.Method = "OPTIONS"
	c.ResponseWriter = &context.Response{
		ResponseWriter: httptest.NewRecorder(),
		Started:        false,
		Status:         0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	MustBeAuthenticated(c)
	Convey("Status should not change", t, func() {
		So(c.ResponseWriter.Status, ShouldEqual, 0)
	})
}
func TestMostBeAuthorizedDeniesAccessIfNotAuthorized(t *testing.T) {
	beego.BConfig.RunMode = "dev"
	r := &http.Request{}
	c := context.NewContext()
	c.Request = r
	c.ResponseWriter = &context.Response{
		ResponseWriter: httptest.NewRecorder(),
		Started:        false,
		Status:         0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	MustBeAuthenticated(c)
	Convey("Status should not change", t, func() {
		So(c.ResponseWriter.Status, ShouldEqual, 400)
	})
	beego.BConfig.RunMode = "test"
}
func TestMostBeAuthorizedTokenMalformed(t *testing.T) {
	beego.BConfig.RunMode = "dev"
	c := context.NewContext()
	rs, _ := http.NewRequest("Get", "/api/v1/", nil)
	rs.Header.Set("Authorization", "Bearer somerandomthing")
	c.Request = rs
	w := httptest.NewRecorder()
	c.ResponseWriter = &context.Response{
		ResponseWriter: w,
		Started:        false,
		Status:         0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	MustBeAuthenticated(c)
	Convey("Status should not change", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 400", func() {
			So(response.Status, ShouldEqual, 400)
		})
		Convey("Error code should be 10013", func() {
			So(response.Code, ShouldEqual, 10013)
		})
		Convey("Error message should inform the user about the admin api ", func() {
			So(response.Message, ShouldEqual, "Token maleformed")
		})
	})
	beego.BConfig.RunMode = "test"
}
func TestMostBeAuthorizedTokenExpired(t *testing.T) {
	beego.BConfig.RunMode = "dev"
	c := context.NewContext()
	rs, _ := http.NewRequest("Get", "/api/v1/", nil)
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	exp, _ := time.Parse(layout, str)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: exp.UTC().Unix(),
		Issuer:    "bc",
	})
	ss, _ := token.SignedString([]byte(beego.AppConfig.String("jwt::key")))

	rs.Header.Set("Authorization", "Bearer "+ss)
	c.Request = rs
	w := httptest.NewRecorder()
	c.ResponseWriter = &context.Response{
		ResponseWriter: w,
		Started:        false,
		Status:         0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	MustBeAuthenticated(c)
	Convey("Status should not change", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 400", func() {
			So(response.Status, ShouldEqual, 400)
		})
		Convey("Error code should be 10013", func() {
			So(response.Code, ShouldEqual, 10013)
		})
		Convey("Error message should inform the user about the admin api ", func() {
			So(response.Message, ShouldEqual, "Token Expired")
		})
	})
	beego.BConfig.RunMode = "test"
}
func TestMostBeAuthorizedCouldNotHandle(t *testing.T) {
	beego.BConfig.RunMode = "dev"
	c := context.NewContext()
	rs, _ := http.NewRequest("Get", "/api/v1/", nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().UTC().Add(time.Hour * time.Duration(1)).Unix(),
		Issuer:    "bc",
	})
	ss, _ := token.SignedString([]byte("abc"))

	rs.Header.Set("Authorization", "Bearer "+ss)
	c.Request = rs
	w := httptest.NewRecorder()
	c.ResponseWriter = &context.Response{
		ResponseWriter: w,
		Started:        false,
		Status:         0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	MustBeAuthenticated(c)
	Convey("Status should not change", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 400", func() {
			So(response.Status, ShouldEqual, 400)
		})
		Convey("Error code should be 10013", func() {
			So(response.Code, ShouldEqual, 10013)
		})
		Convey("Error message should inform the user about the admin api ", func() {
			So(response.Message, ShouldEqual, "Could not handle token")
		})
	})
	beego.BConfig.RunMode = "test"
}

func TestMostBeAuthorizedTokenOK(t *testing.T) {
	beego.BConfig.RunMode = "dev"
	c := context.NewContext()
	rs, _ := http.NewRequest("Get", "/api/v1/", nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).UTC().Unix(),
		Issuer:    "bc",
	})
	ss, _ := token.SignedString([]byte(beego.AppConfig.String("jwt::key")))

	rs.Header.Set("Authorization", "Bearer "+ss)
	c.Request = rs
	w := httptest.NewRecorder()
	c.ResponseWriter = &context.Response{
		ResponseWriter: w,
		Started:        false,
		Status:         0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	MustBeAuthenticated(c)
	Convey("Should Not Write Any Error If Token Not Exp. And Valid", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Status Body Should Be Empty", func() {
			So(len(w.Body.Bytes()), ShouldEqual, 0)
		})
	})
	beego.BConfig.RunMode = "test"
}

func TestPSS(t *testing.T) {
	d, _ := bcrypt.GenerateFromPassword([]byte("artpfmic@123"), bcrypt.DefaultCost)
	fmt.Println(string(d))
}
