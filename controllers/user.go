package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"encoding/json"
	"time"
)

type UserController struct {
	BaseController
}

func (c *UserController) Get() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	users, result := models.GetAllUsers()
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = users
	c.ServeJSON()
}
func (c *UserController) GetByID() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	user := models.User{}
	result := user.GetById(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "User with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}
func (c *UserController) Post() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
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
	createdAt := time.Now()
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
	c.Data["json"] = user
	c.ServeJSON()
}
func (c *UserController) Delete() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	user := models.User{}
	result := user.Delete(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "User with this ID could not be found"
			c.RetError(e)
			return
		} else if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}
func (c *UserController) Update() {
	inData := models.UserEdit{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing UserCreateEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("UserCreateEdit validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed UserCreateEdit:", &inData)
	user := models.User{}
	result := user.Update(c.Ctx.Input.Param(":username"), &inData, time.Now())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}