// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/controllers"
)

func init() {
	// Base route
	beego.Router("/", &controllers.MainController{})
	// Api V1 routes
	ns := beego.NewNamespace("api/v1",
		// Clients Api Endpoints
		beego.NSNamespace("/clients",
			// Get all client
			beego.NSRouter("/", &controllers.ClientController{}, "get:Get"),
			// Insert a client
			beego.NSRouter("/", &controllers.ClientController{}, "post:Post"),
			// Get client by ID
			beego.NSRouter("/:id", &controllers.ClientController{}, "get:GetByID"),
			// Update client
			beego.NSRouter("/:id", &controllers.ClientController{}, "put:Update"),
			// Delete client
			beego.NSRouter("/:id", &controllers.ClientController{}, "delete:Delete"),
		),
	)
	// Add Api v1 namespace to beego.
	beego.AddNamespace(ns)
}
