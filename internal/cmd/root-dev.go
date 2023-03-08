//go:build dev
// +build dev

// 表示只有go build的时候带上tags=dev的时候才会编译这里的内容

package cmd

func init() {
	developmentMode = true
}
