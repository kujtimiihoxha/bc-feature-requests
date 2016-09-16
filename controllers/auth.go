package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/kujtimiihoxha/bc-feature-requests/mail"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"strings"
	"time"
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
			errNotFound.MoreInfo = result.Info
			a.RetError(errNotFound)
			return
		} else if result.Code == models.ErrUnAuthorized {
			errPass := errPermission
			errPass.Message = result.Info
			errPass.MoreInfo = result.Info
			a.RetError(errPass)
			return
		} else if result.Code == models.ErrUserNotVerified {
			errPass := errPermission
			errPass.Message = result.Info
			errPass.MoreInfo = result.Info
			a.RetError(errPass)
			return
		}
	}
	a.Data["json"] = tk
	a.ServeJSON()
}

func MustBeAuthenticated(ctx *context.Context) {
	if ctx.Request.Method == "OPTIONS" {
		return
	}
	if beego.BConfig.RunMode == "test" {
		return
	}
	_, err := ParseToken(ctx)
	if err != nil {
		returnError(ctx, err)
	}
	return

}

// ParseToken parse JWT token in http header.
func ParseToken(ctx *context.Context) (*jwt.Token, *ControllerError) {

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
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				errAuth := errPermission
				errAuth.Message = "Token maleformed"
				return nil, errAuth
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
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
func (base *BaseController) AdminAccessOnly() *ControllerError {
	if beego.BConfig.RunMode == "test" {
		return nil
	}
	tk, err := ParseToken(base.Ctx)
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
func (base *BaseController) NoClientAccessOnly() *ControllerError {
	if beego.BConfig.RunMode == "test" {
		return nil
	}
	tk, err := ParseToken(base.Ctx)
	if err != nil {
		return err
	}
	if tk.Claims.(jwt.MapClaims)["role"] == float64(models.CLIENT) {
		e := errPermission
		e.Message = "This api endpoint is only for admin users"
		return e
	}
	return nil
}

func (c *AuthController) Verify() {
	user := models.User{}
	result := user.Verify(c.Ctx.Input.Param(":id"), time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
		if result.Code == models.ErrNotFound {
			c.RetError(err404)
			return
		}
		if result.Code == models.ErrUserAlreadyVerified {
			err := errInputData
			err.Message = result.Info
			c.RetError(err)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}

func (c *AuthController) Post() {
	inData := models.UserRegister{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing UserCreateEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("UserRegister validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed UserRegister:", &inData)
	createdAt := time.Now().UTC()
	user, result := models.NewUser(&inData, createdAt)
	if result.Code != 0 {
		if result.Code == models.ErrPasswordMismatch || result.Code == models.ErrEmailExists || result.Code == models.ErrUsernameExists {
			er := errInputData
			er.Message = result.Info
			beego.Debug("UserRegister create failed:", er)
			c.RetError(er)
			return
		} else if result.Code == models.ErrSystem {
			er := errSystem
			er.Message = result.Info
			beego.Debug("UserRegister create failed:", er)
			c.RetError(er)
			return
		}
	}
	result = user.Insert()
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	result = mail.Send(user.ID, user.Email)
	if result.Code != 0 {
		user.Delete(user.ID)
		if result.Code == models.ErrEmailNotSent {
			errToSend := errSystem
			errToSend.Message = result.Info
			c.RetError(errToSend)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}

func (c *AuthController) PostClient() {
	inData := models.ClientRegister{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing UserCreateEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("UserRegister validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed UserRegister:", &inData)
	createdAt := time.Now().UTC()
	user, cl, result := models.NewClientUser(&inData, createdAt)
	if result.Code != 0 {
		if result.Code == models.ErrPasswordMismatch || result.Code == models.ErrEmailExists || result.Code == models.ErrUsernameExists {
			er := errInputData
			er.Message = result.Info
			beego.Debug("UserRegister create failed:", er)
			c.RetError(er)
			return
		} else if result.Code == models.ErrSystem {
			er := errSystem
			er.Message = result.Info
			beego.Debug("UserRegister create failed:", er)
			c.RetError(er)
			return
		}
	}
	result = user.InsertClient(cl)
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	result = mail.Send(user.ID, user.Email)
	if result.Code != 0 {
		user.Delete(user.ID)
		cl.Delete(user.ClientID)
		if result.Code == models.ErrEmailNotSent {
			errToSend := errSystem
			errToSend.Message = result.Info
			c.RetError(errToSend)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}
