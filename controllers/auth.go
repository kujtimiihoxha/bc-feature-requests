package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type AuthController struct {
	BaseController
}

func (a *AuthController) Login() {
	userLogin := models.UserLogin{}
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &userLogin)
	if err != nil {
		beego.Debug("Error while parsing ClientCreateEdit:", err)
		a.RetError(errInputData)
		return
	}
	isValid, conErr := a.ValidateInput(userLogin)
	if !isValid {
		beego.Debug("UserLogin validation failed:", conErr)
		a.RetError(conErr)
		return
	}
	user := models.User{}
	tk, result := user.Login(userLogin)
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			a.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			errNotFound := err404
			errNotFound.Message = result.Info
			a.RetError(errNotFound)
			return
		} else if result.Code == models.ErrUnAuthorized {
			errPass := errPermission
			errPass.Message = result.Info
			a.RetError(errPass)
			return
		}
	}
	a.Data["json"] = tk
	a.ServeJSON()
}

func MustBeAuthenticated(ctx *context.Context)  {
	if (beego.BConfig.RunMode == "test"){
		return
	}
	authString := ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		errAuth := errPermission
		errAuth.Message = "Invalid Token"
		returnError(ctx,errAuth)
		return
	}
	tokenString := kv[1]
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwt::key")), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				errAuth := errPermission
				errAuth.Message = "Token maleformed"
				returnError(ctx,errAuth)
				return
			} else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				errAuth := errPermission
				errAuth.Message = "Token Expired"
				returnError(ctx,errAuth)
				return
			} else {
				errAuth := errPermission
				errAuth.Message = "Could not handle token"
				returnError(ctx,errAuth)
			}
		} else {
			// Couldn't handle this token
			errAuth := errPermission
			errAuth.Message = "Could not handle token"
			returnError(ctx,errAuth)
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		errAuth := errPermission
		errAuth.Message = "Invalid Token"
		returnError(ctx,errAuth)
		return
	}
	return

}
// ParseToken parse JWT token in http header.
func  ParseToken(ctx *context.Context) (*jwt.Token, *ControllerError) {

	authString := ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		errAuth := errPermission
		errAuth.Message = "Invalid Token"
		return nil, errAuth
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwt::key")), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				errAuth := errPermission
				errAuth.Message = "Token maleformed"
				return nil, errAuth
			} else if ve.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				errAuth := errPermission
				errAuth.Message = "Token Expired"
				return nil, errAuth
			} else {
				errAuth := errPermission
				errAuth.Message = "Could not handle token"
				return nil, errAuth
			}
		} else {
			// Couldn't handle this token
			errAuth := errPermission
			errAuth.Message = "Could not handle token"
			return nil, errAuth
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		errAuth := errPermission
		errAuth.Message = "Invalid Token"
		return nil, errAuth
	}
	beego.Debug("Token:", token)
	return token, nil
}
func (base *BaseController) AdminAccessOnly() *ControllerError{
	if (beego.BConfig.RunMode == "test"){
		return nil
	}
	tk ,err := ParseToken(base.Ctx)
	if err != nil {
		return err
	}
	if tk.Claims.(jwt.MapClaims)["role"] != float64(models.ADMIN_ROLE) {
		e := errPermission
		e.Message = "This api endpoint is only for admin users"
		return e
	}
	return nil
}