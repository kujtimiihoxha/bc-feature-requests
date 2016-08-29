package controllers

import (
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

type ClientController struct {
	BaseController
}

func (c *ClientController) Get() {
	//TODO add token check.
	clients, result := models.GetAllClients()
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = clients
	c.ServeJSON()
}
func (c *ClientController) GetByID() {
	//TODO add token check.
	client := models.Client{}
	result := client.GetById(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "Client with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}
func (c *ClientController) Post() {
	//TODO add token check.
	inData := models.ClientCreate{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing ClientCreateEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("ClientCreateEdit validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed ClientCreateEdit:", &inData)
	createdAt := time.Now()
	client := models.NewClient(&inData, createdAt)
	result := client.Insert()
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}
func (c *ClientController) Delete() {
	//TODO add token check.
	client := models.Client{}
	result := client.Delete(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "Client with this ID could not be found"
			c.RetError(e)
			return
		} else if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}
func (c *ClientController) Update() {

	//TODO add token check.
	inData := models.ClientEdit{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing ClientCreateEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("ClientCreateEdit validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed ClientCreateEdit:", &inData)
	client := models.Client{}
	result := client.Update(c.Ctx.Input.Param(":id"), &inData, time.Now())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}