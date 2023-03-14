// go:build !dev
//go:build !dev
// +build !dev

package webserver

import (
	"github.com/new-aspect/shiori-practice/internal"
	"io/fs"
)

// assets通常指应用程序中的静态资源，例如图像、CSS和JavaScript文件、模板文件等，
// 它们通常被打包到二进制文件中。使用assets，可以使应用程序更容易地分发、部署和维护。
var assets fs.FS

func init() {
	var err error
	assets, err = fs.Sub(internal.Assets, "view")
	if err != nil {
		panic(err)
	}
}
