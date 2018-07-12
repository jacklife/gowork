package main

import (
	_ "hello/cron"
	_ "hello/routers"
	"github.com/astaxie/beego"
	"path/filepath"
	"hello/util"
)

/***
  beego.LoadAppConfig("ini", filepath.Join(util.GetCfgFilePath(), "app.conf"))
  beego.BConfig.RunMode = beego.DEV
***/
func main() {
	beego.LoadAppConfig("ini", filepath.Join(util.GetCfgFilePath(), "app.conf"))
	beego.BConfig.RunMode = beego.DEV
	beego.SetStaticPath("/swagger",filepath.Join(util.GetCfgFilePath(),"..","swagger"))
	beego.Run()
}
