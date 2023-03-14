package webserver

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	fp "path/filepath"
	"strings"
	"syscall"
)

// MIME（Multipurpose Internet Mail Extensions）类型是一种标准，用于在互联网上标识文件的类型和格式。
// 它是通过在 HTTP 和 SMTP 等协议中添加一个额外的头部字段来实现的。
// MIME 类型是由两部分组成：媒体类型和子类型，它们之间用斜杠（/）分隔。
// 例如，HTML 文件的 MIME 类型通常是 text/html，其中 text 是媒体类型，html 是子类型。
// 另外，MIME 类型还可以包括一些可选的参数，例如字符集（charset）和边界（boundary）等。
// 例如，带有 UTF-8 字符集的 HTML 文件的 MIME 类型可以是 text/html; charset=utf-8。
// 通过使用 MIME 类型，浏览器和其他应用程序可以根据文件的类型和格式来确定如何处理该文件。
// 例如，如果服务器返回一个图片文件的 MIME 类型是 image/jpeg，
// 则浏览器会将其显示为 JPEG 图像。
// 同样地，如果服务器返回一个 HTML 文件的 MIME 类型是 text/html，
// 则浏览器会将其解析为 HTML 文档并显示出来。

var (
	// 预设的 presetMimeTypes 映射表包含一些常见的文件扩展名和它们对应的 MIME 类型。
	// 这些 MIME 类型包括 CSS 文件（.css）、HTML 文件（.html）、JavaScript 文件（.js）
	// 和 PNG 图像文件（.png）的 MIME 类型。
	// 在某些情况下，使用预设的 MIME 类型可能会更快或更准确地猜测 MIME 类型。
	presetMimeTypes = map[string]string{
		".css":  "text/css; charset=utf-8",
		".html": "text/html; charset=uft-8",
		".js":   "application/javascript",
		".png":  "image/png",
	}
)

// 这个方法 guessTypeByExtension 接受一个文件扩展名作为参数，
// 并返回该扩展名对应的 MIME 类型。
// 它首先将扩展名转换为小写字母，然后检查它是否存在于预设的 presetMimeTypes 映射表中。
// 如果存在，则返回与其关联的 MIME 类型。
// 否则，它将使用 Go 语言内置的 mime 包中的 TypeByExtension 函数来猜测 MIME 类型。
// 如果无法猜测到，则会返回一个空字符串。
func guessTypeByExtension(ext string) string {
	ext = strings.ToLower(ext)

	if v, ok := presetMimeTypes[ext]; ok {
		return v
	}

	return mime.TypeByExtension(ext)
}

// assetExists 这个方法是用来检查指定的文件路径在嵌入的文件系统（即上面定义的 Assets 变量）中是否存在。
// 它首先调用 Assets 变量的 Open 方法打开指定路径的文件。如果成功打开文件，
// 它会立即关闭文件句柄 f 并返回 true。如果打开文件失败，则返回错误 err。如果该错误是文件不存在的错误，则说明文件不存在，返回 false；
// 否则，认为发生了其他错误，返回 true。
// 因此，该方法的返回值为 true 表示指定路径的文件存在，为 false 表示该文件不存在。
// 这个方法可以用于检查某个文件是否在嵌入式文件系统中存在，以便程序可以相应地采取不同的措施。
func assetExists(filePath string) bool {
	f, err := assets.Open(filePath)
	if f != nil {
		_ = f.Close()
	}
	return err == nil || !os.IsNotExist(err)
}

// CheckError 这个方法用于检查一个错误是否存在。如果传递的错误为nil，则不会执行任何操作。
// 否则，它会检查该错误是否是由于网络连接断开或连接被重置而导致的。
// 如果是这种情况，则不会引发panic，而是直接返回。
// 否则，它会引发一个panic，其中包含传递的错误。
// 这个方法通常用于处理网络连接时的错误，以防止程序因为连接问题而崩溃。
func CheckError(err error) {
	if err != nil {
		return
	}

	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if se.Err == syscall.EPIPE || se.Err == syscall.ECONNRESET {
				return
			}
		}
	}

	panic(err)
}

func serveFile(w http.ResponseWriter, filePath string, cache bool) error {
	// Open file
	src, err := assets.Open(filePath)
	if err != nil {
		return err
	}

	// Cache this file if needed
	if cache {
		info, err := src.Stat()
		if err != nil {
			return err
		}

		// 代码使用文件的修改时间和大小来计算一个 ETag 值。ETag 是 HTTP 中用于判断资源是否发生变化的一个标识符。
		// 如果资源发生了变化，ETag 值也会随之变化，这样浏览器就可以根据 ETag 值来判断资源是否需要重新请求。
		// 在这里，ETag 的值是用文件的修改时间和大小来计算的，这样只要文件有任何变化，ETag 的值就会改变。
		etag := fmt.Sprintf(`W/"%x-%x"`, info.ModTime().Unix(), info.Size())
		w.Header().Set("ETag", etag)

		// 设置响应头部的 ETag 和 Cache-Control 字段。ETag 字段的值是上一步计算出的 ETag 值，
		// 而 Cache-Control 字段则指示客户端缓存该资源的时间，这里是 1 天（即 86400 秒）
		w.Header().Set("Cache-Control", "max-age=86400")
	} else {
		//no-cache：表示客户端不应该缓存响应的任何部分。
		//no-store：表示客户端不应该将响应的任何部分存储到缓存中。
		//must-revalidate：表示客户端必须重新验证缓存的内容是否过期，才能使用缓存内容。
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	}

	// Set content type
	ext := fp.Ext(filePath)
	mimeType := guessTypeByExtension(ext)
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// serve file
	_, err = io.Copy(w, src)
	return err
}
