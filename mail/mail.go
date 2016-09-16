package mail

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/kujtimiihoxha/bc-feature-requests/models"
	"strconv"
)

var (
	gmail_account          string = beego.AppConfig.String("mail::gmail_account")
	gmail_account_password string = beego.AppConfig.String("mail::gmail_account_password")
	mail_host              string = beego.AppConfig.String("mail::mail_host")
	mail_host_port, err           = beego.AppConfig.Int("mail::mail_host_port")
	email_config           string = `{"username":"` + gmail_account + `","password":"` + gmail_account_password + `","host":"` + mail_host + `","port":` + strconv.Itoa(mail_host_port) + `}`
)

func Send(uuid string, email string) *models.CodeInfo {
	fmt.Println(email_config)
	link := beego.AppConfig.String("server-url") + "/verify/" + uuid
	//
	mail := utils.NewEMail(email_config)
	fmt.Println(mail)
	mail.To = []string{email}
	mail.From = gmail_account
	mail.Subject = "BC Feature Request - Account Activation"
	mail.HTML = "To verify your account, please click on the following link.<br><br><a href=\"" + link +
		"\">" + link + "</a><br><br>Best Regards,<br>BC's team"

	err := mail.Send()

	if err != nil {
		fmt.Println(err)
		beego.Error("SignupController:Post - Unable to send verification email")
		return models.ErrorInfo(models.ErrEmailNotSent, "Unable to send verification email")
	}
	return models.OkInfo("")

}
