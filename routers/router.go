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
	ns := beego.NewNamespace("/api/v1",
		// Clients Api Endpoints
		beego.NSNamespace("/clients",
			// Must be authenticated
			beego.NSBefore(controllers.MustBeAuthenticated),
			// Get all client
			beego.NSRouter("/", &controllers.ClientController{}, "get:Get"),
			// Get client by ID
			beego.NSRouter("/:id", &controllers.ClientController{}, "get:GetByID"),
				// Insert a client
			beego.NSRouter("/", &controllers.ClientController{}, "post:Post"),
			// Update client
			beego.NSRouter("/:id", &controllers.ClientController{}, "put:Update"),
			// Delete client
			beego.NSRouter("/:id", &controllers.ClientController{}, "delete:Delete"),
		),
		beego.NSNamespace("/product-areas",
			// Must be authenticated
			beego.NSBefore(controllers.MustBeAuthenticated),
			// Get all product areas
			beego.NSRouter("/", &controllers.ProductAreaController{}, "get:Get"),
			// Get  product area by ID
			beego.NSRouter("/:id", &controllers.ProductAreaController{}, "get:GetByID"),
			// Insert a  product area
			beego.NSRouter("/", &controllers.ProductAreaController{}, "post:Post"),
			// Update product area
			beego.NSRouter("/:id", &controllers.ProductAreaController{}, "put:Update"),
			// Delete product area
			beego.NSRouter("/:id", &controllers.ProductAreaController{}, "delete:Delete"),
		),
		beego.NSNamespace("/feature-requests",
			beego.NSBefore(controllers.MustBeAuthenticated),
			beego.NSRouter("/", &controllers.FeatureRequestController{}, "post:Post"),
			beego.NSRouter("/", &controllers.FeatureRequestController{}, "get:Get"),
			beego.NSRouter("/:id", &controllers.FeatureRequestController{}, "get:GetByID"),
			beego.NSRouter("/:id/date", &controllers.FeatureRequestController{}, "put:UpdateTargetDate"),
			beego.NSRouter("/:id", &controllers.FeatureRequestController{}, "put:UpdateDetails"),
			beego.NSRouter("/:id/state/:state", &controllers.FeatureRequestController{}, "put:UpdateState"),
			beego.NSRouter("/:id/clients", &controllers.FeatureRequestController{}, "post:AddRemoveClients"),
			beego.NSRouter("/:id/comments", &controllers.FeatureRequestController{}, "post:AddComment"),

		),
		// Users endpoint
		beego.NSNamespace("/users",
			// Must be authenticated
			beego.NSBefore(controllers.MustBeAuthenticated),
			// Update user (only firstname and password)
			beego.NSRouter("/:username", &controllers.UserController{}, "put:Update"),
			// Get all client
			beego.NSRouter("/", &controllers.UserController{}, "get:Get"),
			// Get all employs
			beego.NSRouter("/employs", &controllers.UserController{}, "get:GetEmploys"),
			// Get client by ID
			beego.NSRouter("/:id", &controllers.UserController{}, "get:GetByID"),
			// Delete client
			beego.NSRouter("/:id", &controllers.UserController{}, "delete:Delete"),
		),// Authentication endpoint
		beego.NSNamespace("/auth",
			// Log in
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
			// Verify user
			beego.NSRouter("/verify/:id", &controllers.AuthController{}, "put:Verify"),
			// Register user
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Post"),
		),

	)
	// Add Api v1 namespace to beego.
	beego.AddNamespace(ns)
}
