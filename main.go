package main

import (
	_ "github.com/kujtimiihoxha/bc-feature-requests/routers"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/kujtimiihoxha/bc-feature-requests/controllers"
	"github.com/kujtimiihoxha/bc-feature-requests/db"
	"os"
	"os/signal"
	"syscall"
)

func handleSignals(c chan os.Signal) {
	switch <-c {
	case syscall.SIGINT, syscall.SIGTERM:
		fmt.Println("Shutdown quickly, bye...")
	case syscall.SIGQUIT:
		fmt.Println("Shutdown gracefully, bye...")
		db.Close()
	}
	os.Exit(0)
}

// Application start.
func main() {
	// is graceful shutdown enabled
	graceful, _ := beego.AppConfig.Bool("graceful")
	if !graceful {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go handleSignals(sigs)
	}
	// is logging enabled
	if v, _ := beego.AppConfig.Bool("log::enabled"); v {
		beego.SetLogger(logs.AdapterFile, `{"filename":"`+beego.AppConfig.String("log::path")+`"}`)
	}
	// set logging level to informational if in production mode.
	if beego.BConfig.RunMode == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}
	// set server name.
	beego.BConfig.ServerName = beego.AppConfig.String("server")
	// get the value of cors::allow-credentials
	v, _ := beego.AppConfig.Bool("cors::allow-credentials")
	// set CORS settings
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     beego.AppConfig.Strings("cors::allow-origins"),
		AllowMethods:     beego.AppConfig.Strings("cors::allow-methods"),
		AllowHeaders:     beego.AppConfig.Strings("cors::allow-headers"),
		ExposeHeaders:    beego.AppConfig.Strings("cors::expose-headers"),
		AllowCredentials: v,
	}))
	// specify error controller
	beego.ErrorController(&controllers.ErrorController{})
	// connect to the db (if no connection panic)
	db.Connect()
	// run the server
	beego.Run()
}
