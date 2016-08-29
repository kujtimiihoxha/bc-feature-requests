package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"github.com/asaskevich/govalidator"
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


// ParseToken parse JWT token in http header.
func (base *BaseController) ParseToken() (*jwt.Token, *ControllerError) {
	authString := base.Ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil, errInputData
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				return nil, errInputData
			} else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, errExpired
			} else {
				// Couldn't handle this token
				return nil, errInputData
			}
		} else {
			// Couldn't handle this token
			return nil, errInputData
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		return nil, errInputData
	}
	beego.Debug("Token:", token)
	return token, nil
}
