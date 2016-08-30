package controllers

import (
	"github.com/astaxie/beego"
	"github.com/asaskevich/govalidator"
	"github.com/astaxie/beego/context"
	"encoding/json"
)

// BaseController.
// The base controller used from all controllers (besides MainController).
type BaseController struct {
	beego.Controller
}
// RetError.
// Returns errors from other controllers.
func (base *BaseController) RetError(e *ControllerError) {
	if mode := beego.AppConfig.String("runmode"); mode == "prod" {
		e.DevInfo = ""
	}
	base.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	base.Ctx.ResponseWriter.WriteHeader(e.Status)
	base.Data["json"] = e
	base.ServeJSON()
	base.StopRun()
}
func (base *BaseController) ValidateInput(i interface{}) (bool, *ControllerError) {
	isValid,err := govalidator.ValidateStruct(i)
	if err != nil {
		controllerError := errInputDataValidation;
		controllerError.MoreInfo = err.Error()
		return false, controllerError
	}
	if !isValid {
		return false, errInputDataValidation
	}
	return isValid,nil
}

func returnError(ctx *context.Context, e *ControllerError) {
	if mode := beego.AppConfig.String("runmode"); mode == "prod" {
		e.DevInfo = ""
	}
	ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	ctx.ResponseWriter.WriteHeader(e.Status)
	d,_ := json.Marshal(e)
	ctx.ResponseWriter.Write(d)
}