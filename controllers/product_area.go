package controllers

import (
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

type ProductAreaController struct {
	BaseController
}

func (c *ProductAreaController) Get() {
	product_areas, result := models.GetAllProductAreas()
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = product_areas
	c.ServeJSON()
}
func (c *ProductAreaController) GetByID() {
	product_area := models.ProductArea{}
	result := product_area.GetById(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "ProductArea with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = product_area
	c.ServeJSON()
}
func (c *ProductAreaController) Post() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	inData := models.ProductAreaCreate{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing ProductAreaEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("ProductAreaEdit validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed ProductAreaEdit:", &inData)
	createdAt := time.Now().UTC()
	product_area := models.NewProductArea(&inData, createdAt)
	result := product_area.Insert()
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}
	}
	c.Data["json"] = product_area
	c.ServeJSON()
}
func (c *ProductAreaController) Delete() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	product_area := models.ProductArea{}
	result := product_area.Delete(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "ProductArea with this ID could not be found"
			c.RetError(e)
			return
		} else if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		}
	}
	c.Data["json"] = product_area
	c.ServeJSON()
}
func (c *ProductAreaController) Update() {
	res := c.AdminAccessOnly()
	if res != nil {
		beego.Debug("No Access", res)
		c.RetError(res)
		return
	}
	inData := models.ProductAreaEdit{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing ProductAreaEdit:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("ProductAreaEdit validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed ProductAreaeEdit:", &inData)
	product_area := models.ProductArea{}
	result := product_area.Update(c.Ctx.Input.Param(":id"), &inData, time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "ProductArea with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = product_area
	c.ServeJSON()
}