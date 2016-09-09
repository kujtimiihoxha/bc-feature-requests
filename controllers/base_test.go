package controllers

import (
	_ "github.com/kujtimiihoxha/bc-feature-requests/tests"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/astaxie/beego"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"net/http/httptest"
	"github.com/astaxie/beego/context"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
)

func TestGetUser(t *testing.T) {
	user := models.User{
		BaseModel:models.BaseModel{
			ID:"585aca07-6d3c-43ba-97d4-8fb4cb27e024",
		},
		Username:"employ",
		Email:"employ@gmail.com",
		Role:2,
		FirstName:"Employ",
		LastName:"Name",
		Password:"$2a$10$Rtt0sfArkW1gCLeiW5AUbu6VgRNtzzYRKPmD5xmK/JhAyw4VA8Ipq",
	}
	c := context.NewContext()
	rs, _ := http.NewRequest("Get", "/api/v1/",nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt:time.Now().Add(time.Hour).UTC().Unix(),
			Issuer:    "bc",
		},
		Username: user.Username,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Role: user.Role,
		ID: user.ID,
	})
	ss, _ := token.SignedString([]byte(beego.AppConfig.String("jwt::key")))

	rs.Header.Set("Authorization","Bearer "+ss)
	c.Request = rs
	w := httptest.NewRecorder()
	c.ResponseWriter = &context.Response{
		ResponseWriter:w,
		Started:false,
		Status:0,
	}
	c.Output.Reset(c)
	c.Input.Reset(c)
	base := BaseController{
		Controller:beego.Controller{
			Ctx:c,
		},
	}
	Convey("Test if login succesfull\n", t, func() {
		So(base.User().ID,ShouldEqual,user.ID)
	})
}