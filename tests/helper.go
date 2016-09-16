package tests

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"runtime"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "."+string(filepath.Separator))))
	beego.TestBeegoInit(appPath)
}
