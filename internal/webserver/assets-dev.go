//go:build dev
// +build dev

package webserver

import "net/http"

// 这是做什么的？
var assets = http.Dir("internal/view")

func init() {
	developmentMode = true
}
