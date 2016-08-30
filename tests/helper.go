package tests

import (
	"github.com/astaxie/beego"
	"runtime"
	"path/filepath"
)


func init() {
	_, file, _, _ := runtime.Caller(1)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "." + string(filepath.Separator))))
	beego.TestBeegoInit(appPath)
}

