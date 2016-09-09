package controllers

import (
	"net/http/httptest"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"encoding/json"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"github.com/astaxie/beego"
)


func ValidationFail(w *httptest.ResponseRecorder, t *testing.T, messge string) {
	Convey("Test id insert fails if data not valid\n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Code should be 10014", func() {
			So(response.Code, ShouldEqual,10014)
		})
		Convey("MoreInfo should tell the user what he forgot", func() {
			So(response.MoreInfo, ShouldEqual,messge)
		})
	})
}

func FailNoData(w *httptest.ResponseRecorder, t *testing.T) {
	Convey("Test if insert fails when no data \n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Code should be 10001", func() {
			So(response.Code, ShouldEqual,10001)
		})
		Convey("Message should specify the Input data error", func() {
			So(response.Message, ShouldEqual,ErrInputData)
		})
	})
}
func SystemError(w *httptest.ResponseRecorder, t *testing.T) {
	Convey("Test if system error is returned \n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, 500)
		})
		Convey("The Result Code should be 10011", func() {
			So(response.Code, ShouldEqual,10011)
		})
		Convey("Message should specify the system erorr message", func() {
			So(response.Message, ShouldEqual,ErrSystem)
		})
	})
}
func DataBaseError(w *httptest.ResponseRecorder, t *testing.T) {
	Convey("Test if database error is returned \n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, 500)
		})
		Convey("The Result Code should be 10002", func() {
			So(response.Code, ShouldEqual,10002)
		})
		Convey("Message should specify the error as a database error", func() {
			So(response.Message, ShouldEqual,ErrDatabase)
		})
	})

}
func Insert(w *httptest.ResponseRecorder, t *testing.T,id string) {
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
func Update(w *httptest.ResponseRecorder, t *testing.T,inData models.ClientEdit) {
	Convey("Test if id is set after success update\n", t, func() {
		response := models.Client{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Response client should have the name of the edit data", func() {
			So(response.Name, ShouldEqual,inData.Name)
		})
		Convey("Response client should have the description of the edit data", func() {
			So(response.Description, ShouldEqual,inData.Description)
		})
	})
}
func NotFound(w *httptest.ResponseRecorder, t *testing.T, moreInfo string)  {
	Convey("Test if not found is returned \n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(),&response)
		Convey("Status Code Should Be 404", func() {
			So(w.Code, ShouldEqual, 404)
		})
		Convey("The Result Code should be 404", func() {
			So(response.Code, ShouldEqual,404)
		})
		Convey("More info should specify the needed message", func() {
			So(response.MoreInfo, ShouldEqual,moreInfo)
		})
	})
}
func NotAuthorized(w *httptest.ResponseRecorder, t *testing.T)  {
	Convey("Test access is denied for non admin users\n", t, func() {
		response := ControllerError{}
		json.Unmarshal(w.Body.Bytes(), &response)
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("Error code should be 10013", func() {
			So(response.Code, ShouldEqual, 10013)
		})
		Convey("Error message should inform the user about the admin api ", func() {
			So(response.Message, ShouldEqual, "This api endpoint is only for admin users")
		})
	})
	beego.BConfig.RunMode = "test"
}