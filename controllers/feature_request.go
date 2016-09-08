package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"encoding/json"
	"time"
	"strconv"
	"fmt"
)

type FeatureRequestController struct {
	BaseController
}

func (fr *FeatureRequestController) Post() {
	inData := models.FeatureRequestCreate{}
	err  := json.Unmarshal(fr.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing ClientCreateEdit:", err)
		fr.RetError(errInputData)
		return
	}
	isValid, conErr := fr.ValidateInput(inData)
	if !isValid {
		beego.Debug("ClientCreateEdit validation failed:", conErr)
		fr.RetError(conErr)
		return
	}
	beego.Debug("Parsed ClientCreateEdit:", &inData)
	createdAt := time.Now().UTC()
	featureRequest := models.NewFeatureRequest(&inData, createdAt,fr.User().ID)
	result := featureRequest.Insert()
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			fr.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			fr.RetError(errSystem)
			return
		}
	}
	fr.Data["json"] = featureRequest
	fr.ServeJSON()
}
func (fr *FeatureRequestController) AddComment() {
	inData := models.FeatureRequestAddComment{}
	err  := json.Unmarshal(fr.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing FeatureRequestAddComment:", err)
		fr.RetError(errInputData)
		return
	}
	isValid, conErr := fr.ValidateInput(inData)
	if !isValid {
		beego.Debug("FeatureRequestAddComment validation failed:", conErr)
		fr.RetError(conErr)
		return
	}
	beego.Debug("Parsed FeatureRequestAddComment:", &inData)
	featureRequest := models.FeatureRequest{}
	result := featureRequest.AddComment(fr.Ctx.Input.Param(":id"),&inData)
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			fr.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			fr.RetError(errSystem)
			return
		}
	}
	fr.Data["json"] = featureRequest
	fr.ServeJSON()
}


func (c *FeatureRequestController) GetByID() {
	featureRequest := models.FeatureRequest{}
	result := featureRequest.GetById(c.Ctx.Input.Param(":id"))
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		} else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "FeatureRequest with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = featureRequest
	c.ServeJSON()
}
func (fr *FeatureRequestController) Get() {
	filter := fr.ParseFilter()

	feature_requests,result := models.GetFeatureRequestByFilterSort(&filter)
	if result.Code != 0 {
		if result.Code == models.ErrSystem {
			fr.RetError(errSystem)
			return
		}
	}
	fr.Data["json"] = feature_requests
	fr.ServeJSON()
}
func (fr *FeatureRequestController) ParseFilter() models.FeatureRequestFilter {
	cl,_ := fr.GetInt("closed");
	skip,_ := fr.GetInt("skip");
	get,_ := fr.GetInt("get");
	return models.FeatureRequestFilter{
		Client:fr.GetString("client"),
		Closed:cl,
		Employ:fr.GetString("employ"),
		ProductArea:fr.GetString("product_area"),
		ClientPriorityDir:fr.GetString("priority_dir"),
		FeatureRequestSort:models.FeatureRequestSort{
			Dir:fr.GetString("dir"),
			Field:fr.GetString("field"),
		},
		FeatureRequestPagination: models.FeatureRequestPagination{
			Skip: skip,
			Get: get,
		},
	}
}

func (c *FeatureRequestController) UpdateTargetDate() {
	inData := models.FeatureRequestEditTargetDate{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing FeatureRequestEditTargetDate:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("FeatureRequestEditTargetDate validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed FeatureRequestEditTargetDate:", &inData)
	client := models.FeatureRequest{}
	result := client.UpdateTargetDate(c.Ctx.Input.Param(":id"),c.User().ID,c.User().Username, &inData, time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "Feature request with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}

func (c *FeatureRequestController) AddRemoveClients() {
	inData := models.FeatureRequestAddRemoveClients{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing FeatureRequestAddRemoveClients:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("FeatureRequestAddRemoveClients validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed FeatureRequestAddRemoveClients:", &inData)
	client := models.FeatureRequest{}
	result := client.AddRemoveClients(c.Ctx.Input.Param(":id"),c.User().ID,c.User().Username, &inData, time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "Feature request with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}

func (c *FeatureRequestController) UpdateDetails() {
	inData := models.FeatureRequestEditDetails{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &inData)
	if err != nil {
		beego.Debug("Error while parsing FeatureRequestEditDetails:", err)
		c.RetError(errInputData)
		return
	}
	isValid, conErr := c.ValidateInput(inData)
	if !isValid {
		beego.Debug("FeatureRequestEditDetails validation failed:", conErr)
		c.RetError(conErr)
		return
	}
	beego.Debug("Parsed FeatureRequestEditTargetDate:", &inData)
	client := models.FeatureRequest{}
	result := client.UpdateDetails(c.Ctx.Input.Param(":id"),c.User().ID,c.User().Username, &inData, time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "Feature request with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}

func (c *FeatureRequestController) UpdateState() {
	inData,err := strconv.ParseInt(c.Ctx.Input.Param(":state"),10,0)
	fmt.Println(inData,err)
	if err!= nil || (inData != 1 && inData != 2 ) {
		controllerError := errInputDataValidation;
		controllerError.MoreInfo = "State must be either 1 or 2"
		beego.Debug("State validation failed:", controllerError)
		c.RetError(controllerError)
		return
	}
	beego.Debug("Parsed inData:", inData)
	client := models.FeatureRequest{}
	state := false;
	if inData == 1 {
		state = true
	}
	result := client.UpdateState(c.Ctx.Input.Param(":id"),c.User().ID,c.User().Username, state, time.Now().UTC())
	if result.Code != 0 {
		if result.Code == models.ErrDatabase {
			c.RetError(errDatabase)
			return
		} else if result.Code == models.ErrSystem {
			c.RetError(errSystem)
			return
		}else if result.Code == models.ErrNotFound {
			e := err404
			e.MoreInfo = "Feature request with this ID could not be found"
			c.RetError(e)
			return
		}
	}
	c.Data["json"] = client
	c.ServeJSON()
}