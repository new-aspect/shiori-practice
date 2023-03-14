package webserver

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path"
	fp "path/filepath"
	"strings"
)

// serveJsFile is handler for GET /js/*filepath
func (h *handler) serveJsFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsFilePath := ps.ByName("filepath")
	// todo 这里好像是利用path.Join拼接路径
	jsFilePath = path.Join("js", jsFilePath)
	// 这里是利用path.split区分出路径和文件名, 这样命名也太优雅了
	jsDir, jsName := path.Split(jsFilePath)

	if developmentMode && fp.Ext(jsName) == ".js" && strings.HasSuffix(jsName, ".min.js") {
		jsName = strings.TrimSuffix(jsName, ".min.js") + ".js"
		tmpPath := path.Join(jsDir, jsName)
		if assetExists(tmpPath) {
			jsFilePath = tmpPath
		}
	}

	//todo 这个serveFile需要实现
	err := serveFile(w, jsFilePath, true)
	CheckError(err)
}

func (h *handler) serveFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// todo 这里没有完成
}
