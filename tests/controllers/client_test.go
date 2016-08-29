package tests_controllers

import (
	"testing"
	r "github.com/dancannon/gorethink"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	_ "github.com/kujtimiihoxha/bc-feature-requests/routers"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/astaxie/beego"
	"runtime"
	"path/filepath"
	"net/http"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"github.com/kujtimiihoxha/bc-feature-requests/controllers"
	"github.com/iris-contrib/errors"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"time"
	"reflect"
)
func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../" + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}
var clients_table = "clients"

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
	b,_ := json.Marshal(insertClientNoName)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test id insert fails if data not valid\n", t, func() {
		response := controllers.ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Code should be 10014", func() {
			So(response.Code, ShouldEqual,10014)
		})
		Convey("MoreInfo should tell the user what he forgot (in this case Name)", func() {
			So(response.MoreInfo, ShouldEqual,"Name: non zero value required;")
		})
	})
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
	Convey("Test if insert fails when no data \n", t, func() {
		response := controllers.ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Code should be 10001", func() {
			So(response.Code, ShouldEqual,10001)
		})
		Convey("Message should specify the Input data error", func() {
			So(response.Message, ShouldEqual,controllers.ErrInputData)
		})
	})
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
	b,_ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Insert(r.MockAnything())).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if system error is returned \n", t, func() {
		response := controllers.ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, 500)
		})
		Convey("The Result Code should be 10011", func() {
			So(response.Code, ShouldEqual,10011)
		})
		Convey("Message should specify the system erorr message", func() {
			So(response.Message, ShouldEqual,controllers.ErrSystem)
		})
	})
}

// Test the response when there is any database error.
// Should Return :
// Status : 500
// Code : 10002
func TestDatabaseError(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b,_ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Insert(r.MockAnything())).Once().Return(r.WriteResponse{ Errors:1,FirstError:"Error from DB"}, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if database error is returned \n", t, func() {
		response := controllers.ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, 500)
		})
		Convey("The Result Code should be 10002", func() {
			So(response.Code, ShouldEqual,10002)
		})
		Convey("Message should specify the error as a database error", func() {
			So(response.Message, ShouldEqual,controllers.ErrDatabase)
		})
	})
}

// Test the response when there is any database error.
// Should Return :
// Status : 200
// Client : With ID set
func TestInsertSuccess(t *testing.T) {
	id := "585aca07-6d3c-43ba-97d4-8fb4cb27e024"
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b,_ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table).Insert(r.MockAnything())).Once().Return(r.WriteResponse{GeneratedKeys:[]string{id}}, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if id is set after success insert\n", t, func() {
		response := models.Client{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Response client should have the generated id set", func() {
			So(response.ID, ShouldEqual,id)
		})
	})
}


/*-------------------------------
	Test Get
--------------------------------*/

// Test the response when there is any system error.
// Should Return :
// Status : 500
// Code : 10011
func TestGetAllSystemError(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b,_ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("GET", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	mock.On(r.Table(clients_table)).Once().Return(nil, errors.New("Test error"))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if system error is returned \n", t, func() {
		response := controllers.ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, 500)
		})
		Convey("The Result Code should be 10011", func() {
			So(response.Code, ShouldEqual,10011)
		})
		Convey("Message should specify the system erorr message", func() {
			So(response.Message, ShouldEqual,controllers.ErrSystem)
		})
	})
}


// Test the response when successful
// Should Return :
// Status : 200
// Clients: array of clients in the DB.
func TestGetAllSuccess(t *testing.T) {
	insertClient := models.ClientCreate{
		Description:"Test",
		Name:"Test",
	}
	b,_ := json.Marshal(insertClient)
	rs, _ := http.NewRequest("GET", "/api/v1/clients", bytes.NewReader(b))
	mock := r.NewMock()
	db.SetTestSession(mock)
	tm := time.Now()
	dbClients :=[]models.Client{
		{
			BaseModel:models.BaseModel{
				ID: "585aca07-6d3c-43ba-97d4-8fb4cb27e024",
				CreatedAt:&tm,
			},
			Name:"Test",
			Description:"Desc",

		}}
	mock.On(r.Table(clients_table)).Once().Return(dbClients,nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if system error is returned \n", t, func() {
		clients := []models.Client{}
		json.Unmarshal(w.Body.Bytes(),&clients)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Should return an array of existing clients", func() {
			So(len(clients), ShouldEqual,1)
		})
		Convey("Data should be equal to the data in the DB", func() {
			So(reflect.DeepEqual(clients,dbClients), ShouldEqual,true)
		})
	})
}