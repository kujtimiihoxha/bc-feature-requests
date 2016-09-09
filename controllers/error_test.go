package controllers
import (
	"testing"
	"github.com/astaxie/beego"
	_ "github.com/kujtimiihoxha/bc-feature-requests/tests"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego/context"
)

func TestErrorController(t *testing.T) {
	c := context.NewContext()
	rs, _ := http.NewRequest("Get", "/api/v1/",nil)
	c.Request = rs
	w := httptest.NewRecorder()
	c.ResponseWriter = &context.Response{
		ResponseWriter:w,
		Started:false,
		Status:0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	errorC := ErrorController{
		BaseController:BaseController{
			Controller: beego.Controller{
				Ctx:c,
				Data:map[interface {}]interface {}{},
			},
		},
	}
	Convey("Test if error controller works \n", t, func() {
		Convey("Should stop runing user request", func() {
			ShouldPanic(func() {
				errorC.Error404()
			})
		})
		Convey("Should have status set",func() {
			So(errorC.Ctx.ResponseWriter.Status, ShouldEqual, 404)
		})
	})
}