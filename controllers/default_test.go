package controllers

import (
	"testing"
	"github.com/astaxie/beego"
	_ "github.com/kujtimiihoxha/bc-feature-requests/tests"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"reflect"
)
// Init client routes  for testing.
func init() {
	// Base route
	beego.Router("/", &MainController{})
}

func TestDefaultRoute(t *testing.T) {
	rs, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, rs)
	Convey("Test if main route returnes api helper \n", t, func() {
		apiHelp :=  map[string]string{}
		json.Unmarshal(w.Body.Bytes(), &apiHelp)
		So(reflect.DeepEqual(apiHelp,API_HELP), ShouldEqual, true)
	})
}