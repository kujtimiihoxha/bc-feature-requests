package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"time"
)

type ClientController struct {
	BaseController
}

func (c *ClientController) Get() {
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
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
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
	createdAt := time.Now().UTC()
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
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
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
		} else if result.Code == models.ErrRecordHasConnections {
			err := errInputData
			err.Message = result.Info
			c.RetError(err)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}
func (c *ClientController) Update() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
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
	result := client.Update(c.Ctx.Input.Param(":id"), &inData, time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
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
