package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"time"
)

type UserController struct {
	BaseController
}

func (c *UserController) Get() {
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
func (c *UserController) GetNotifications() {
	user := models.User{}
	notifications, result := user.GetNotifications(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = notifications
	c.ServeJSON()
}
func (c *UserController) SetNotificationsViewed() {
	user := models.User{}
	result := user.SetNotificationsViewed(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)

		}

		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		}
	}
	c.Data["json"] = user
	c.ServeJSON()
}
func (c *UserController) GetEmploys() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	users, result := models.GetAllEmploys()
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
	result := user.Update(c.Ctx.Input.Param(":username"), &inData, time.Now().UTC())
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
