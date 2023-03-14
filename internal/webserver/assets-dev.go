//go:build dev
// +build dev

package webserver

import "net/http"

// 注意，assets-dev表示开发环境处理，assets-prod表示生产环境处理

// http.Dir 函数返回一个实现了 http.FileSystem 接口的类型，
// 用于表示一个本地文件系统目录的路径。
// 因此，这行代码的作用是将 "internal/view" 目录
// 的路径保存在 assets 变量中，以便在后续的程序中使用。
var assets = http.Dir("internal/view")

func init() {
	developmentMode = true
}
