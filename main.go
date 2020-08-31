package main

import (
	_ "goku/routers"
	"goku/types"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

func main() {
	logLevel := beego.AppConfig.String("loglevel")
	switch logLevel {
	case "trace":
		logs.SetLevel(logs.LevelTrace)
		break
	case "debug":
		logs.SetLevel(logs.LevelDebug)
		break
	case "info":
		logs.SetLevel(logs.LevelInfo)
		break
	}
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	types.InitMysql()
	beego.Run()
}
