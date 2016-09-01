package controllers

import (
	"testing"
	r "github.com/dancannon/gorethink"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	_ "github.com/kujtimiihoxha/bc-feature-requests/tests"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/astaxie/beego"
	"net/http"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"github.com/iris-contrib/errors"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"time"
	"reflect"
	"github.com/dgrijalva/jwt-go"
)

var clients_table = "clients"
// Init client routes  for testing.
func init() {
	// Base route
	beego.Router("/", &MainController{})
	// Api V1 routes
	ns := beego.NewNamespace("/api/v1",
		// Clients Api Endpoints
		beego.NSNamespace("/clients",
			// Must be authenticated
			beego.NSBefore(MustBeAuthenticated),
			// Get all client
			beego.NSRouter("/", &ClientController{}, "get:Get"),
			// Get client by ID
			beego.NSRouter("/:id", &ClientController{}, "get:GetByID"),
			// Insert a client
			beego.NSRouter("/", &ClientController{}, "post:Post"),
			// Update client
			beego.NSRouter("/:id", &ClientController{}, "put:Update"),
			// Delete client
			beego.NSRouter("/:id", &ClientController{}, "delete:Delete"),
		),
	)
	// Add Api v1 namespace to beego.
	beego.AddNamespace(ns)
}

/*-------------------------------
	Test Insert
--------------------------------*/

// Tests if the insert method fails
//  if the data send is not valid.
// Should Return :
// Status : 400
// Code : 10014
func TestInsertFailValidation(t *testing.T) {
	insertClientNoName := models.ClientCreate{
		Description:"Test",
	}
	b, _ := json.Marshal(insertClientNoName)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	ValidationFail(w, t, "Name: non zero value required;")
}

// Tests if the insert method fails
//  if no data is sent.
// Should Return :
// Status : 400
// Code : 10001
func TestInsertFailWhenNoData(t *testing.T) {
	rs, _ := http.NewRequest("POST", "/api/v1/clients", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	FailNoData(w, t)
}

// Test the response when there is any system error.
// Should Return :
// Status : 500
// Code : 10011
func TestInsertSystemError(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b, _ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Insert(r.MockAnything())).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	SystemError(w, t)
}

// Test the response when there is any database error.
// Should Return :
// Status : 500
// Code : 10002
func TestInsertDatabaseError(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b, _ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Insert(r.MockAnything())).Once().Return(r.WriteResponse{Errors:1, FirstError:"Error from DB"}, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	DataBaseError(w, t)
}

// Test the response when there is any database error.
// Should Return :
// Status : 200
// Client : With ID set
func TestInsertSuccess(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	b, _ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Insert(r.MockAnything())).Once().Return(r.WriteResponse{GeneratedKeys:[]string{id}}, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Insert(w, t, id)
}

// Test the response when user has no privilege.
// Should Return :
// Status : 400
// Code : 10013
func TestInsertOnlyAdmin(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b, _ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	w := httptest.NewRecorder()
	claims := &models.Claims{
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "bc",
		},
		Username: "username",
		Role: 2,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(beego.AppConfig.String("jwt::key")))
	rs.Header.Add("Authorization","Bearer "+ss)
	beego.BConfig.RunMode = "dev"
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	NotAuthorized(w,t)
}


/*-------------------------------
	Test Get
--------------------------------*/

// Test the response when there is any system error.
// Should Return :
// Status : 500
// Code : 10011
func TestGetAllSystemError(t *testing.T) {
	rs, _ := http.NewRequest("GET", "/api/v1/clients", nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table)).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	SystemError(w, t)
}


// Test the response when successful
// Should Return :
// Status : 200
// Clients: array of clients in the DB.
func TestGetAllSuccess(t *testing.T) {
	rs, _ := http.NewRequest("GET", "/api/v1/clients", nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	tm := time.Now()
	dbClients := []models.Client{
		{
			BaseModel:models.BaseModel{
				ID: "585aca07-6d3c-43ba-97d4-8fb4cb27e024",
				CreatedAt:&tm,
			},
			Name:"Test",
			Description:"Desc",

		}}
	mock.On(r.Table(clients_table)).Once().Return(dbClients, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if system error is returned \n", t, func() {
		clients := []models.Client{}
		json.Unmarshal(w.Body.Bytes(), &clients)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Should return an array of existing clients", func() {
			So(len(clients), ShouldEqual, 1)
		})
		Convey("Data should be equal to the data in the DB", func() {
			So(reflect.DeepEqual(clients, dbClients), ShouldEqual, true)
		})
	})
}

/*-------------------------------
	Test GetById
--------------------------------*/

// Test the response when there is any system error.
// Should Return :
// Status : 500
// Code : 10011
func TestGetByIDSystemError(t *testing.T) {
	rs, _ := http.NewRequest("GET", "/api/v1/clients/585aca07-6d3c-43ba-97d4-8fb4cb27e024", nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get("585aca07-6d3c-43ba-97d4-8fb4cb27e024")).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	SystemError(w, t)
}

// Test the response when the record with the id could not be found.
// Should Return :
// Status : 404
// Code : 404
func TestGetByIDNotFound(t *testing.T) {
	rs, _ := http.NewRequest("GET", "/api/v1/clients/585aca07-6d3c-43ba-97d4-8fb4cb27e024", nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get("585aca07-6d3c-43ba-97d4-8fb4cb27e024")).Once().Return(nil, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	NotFound(w, t)
}


// Test the response when successful
// Should Return :
// Status : 200
// Clients: model of clients in the DB.
func TestGetByIDSuccess(t *testing.T) {
	client := models.Client{
		Name:"Test",
		Description:"test",
	}
	rs, _ := http.NewRequest("GET", "/api/v1/clients/585aca07-6d3c-43ba-97d4-8fb4cb27e024", nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get("585aca07-6d3c-43ba-97d4-8fb4cb27e024")).Once().Return(client, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if get by id success \n", t, func() {
		cl := models.Client{}
		json.Unmarshal(w.Body.Bytes(), &cl)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Data should be equal to the data in the DB", func() {
			So(reflect.DeepEqual(client, cl), ShouldEqual, true)
		})
	})
}


/*-------------------------------
	Test Delete
--------------------------------*/

// Test the response when there is any system error.
// Should Return :
// Status : 500
// Code : 10011
func TestDeleteSystemError(t *testing.T) {
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	rs, _ := http.NewRequest("DELETE", "/api/v1/clients/" + id, nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get(id)).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	SystemError(w, t)
}

// Test the response when the record with the id could not be found.
// Should Return :
// Status : 404
// Code : 404
func TestDeleteNotFound(t *testing.T) {
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	rs, _ := http.NewRequest("DELETE", "/api/v1/clients/" + id, nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get(id)).Once().Return(nil, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	NotFound(w, t)
}


// Test the response when there is any database error.
// Should Return :
// Status : 500
// Code : 10002
func TestDeleteDatabaseError(t *testing.T) {
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	rs, _ := http.NewRequest("DELETE", "/api/v1/clients/" + id, nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get(id)).Twice().Return(models.Client{
		BaseModel: models.BaseModel{
			ID:"585aca07-6d3c-43ba-97d4-8fb4cb27e024",
		},
	}, nil)
	mock.On(r.Table(clients_table).Get(id).Delete()).Once().Return(r.WriteResponse{Errors:1, FirstError:"Error from DB"}, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	DataBaseError(w, t)
}
// Test the response when successful
// Should Return :
// Status : 200
// Clients: array of clients in the DB.
func TestDeleteSuccess(t *testing.T) {
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	dbClient := models.Client{
		BaseModel:models.BaseModel{
			ID: id,
		},
		Name:"Test",
		Description:"Desc",

	}
	rs, _ := http.NewRequest("DELETE", "/api/v1/clients/" + id, nil)
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Get(id)).Once().Return(dbClient, nil)
	mock.On(r.Table(clients_table).Get(id).Delete()).Once().Return(r.WriteResponse{Deleted:1}, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if response is success \n", t, func() {
		client := models.Client{}
		json.Unmarshal(w.Body.Bytes(), &client)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Data should be equal to the data in the DB", func() {
			So(reflect.DeepEqual(client, dbClient), ShouldEqual, true)
		})
	})
}

// Test the response when user has no privilege.
// Should Return :
// Status : 400
// Code : 10013
func TestDeleteOnlyAdmin(t *testing.T) 	{
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	rs, _ := http.NewRequest("DELETE", "/api/v1/clients/" + id, nil)
	w := httptest.NewRecorder()
	claims := &models.Claims{
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "bc",
		},
		Username: "username",
		Role: 2,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(beego.AppConfig.String("jwt::key")))
	rs.Header.Add("Authorization","Bearer "+ss)
	beego.BConfig.RunMode = "dev"
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	NotAuthorized(w,t)
}
