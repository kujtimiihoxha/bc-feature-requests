package controllers

import (
	"github.com/astaxie/beego"
)
// MainController.
// The controller of the main route.
// Route: /
// Methods allowed: GET
type MainController struct {
	beego.Controller
}
// Server base url
var baseUrl = beego.AppConfig.String("server-url")

// Api help map.
// The api help map contains useful api endpoints that van be accessed.
var API_HELP map[string]string = map[string]string{
	"clients":baseUrl+"/api/v1/clients",
}
// Get.
// Get the api help json.
// Route: /
// Method: GET
func (c *MainController) Get() {
	c.Data["json"]= API_HELP
	c.ServeJSON()
}

// No other method is allowed in the route: /
