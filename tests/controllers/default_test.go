package tests_controllers

import (
	"testing"
	"github.com/astaxie/beego"
	"runtime"
	"path/filepath"
	_ "github.com/kujtimiihoxha/bc-feature-requests/routers"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"github.com/kujtimiihoxha/bc-feature-requests/controllers"
	"reflect"
)
func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../" + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}
// Test the main controller.
// Route : /
// Method Get.
func TestMainController(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test Base Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			response := map[string]string{}
			json.Unmarshal(w.Body.Bytes(),&response)
			So(reflect.DeepEqual(response, controllers.API_HELP), ShouldEqual,true)
		})
	})
}