package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"encoding/json"
	"time"
	"github.com/dgrijalva/jwt-go"
)

type FeatureRequestController struct {
	BaseController
}

func (fr *FeatureRequestController) Post() {
	token,terr := ParseToken(fr.Ctx);
	if terr != nil {
		fr.RetError(terr)
	}
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
	createdAt := time.Now()
	featureRequest := models.NewFeatureRequest(&inData, createdAt,token.Claims.(jwt.MapClaims)["id"].(string))
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