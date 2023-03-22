package webserver

import (
	"github.com/julienschmidt/httprouter"
	"log"
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

// serveFile is handler for general file request
func (h *handler) serveFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// filePath = css/archive.css
	rootPath := strings.Trim(h.RootPath, "/")
	urlPath := strings.Trim(r.URL.Path, "/")
	filePath := strings.TrimPrefix(urlPath, rootPath)
	filePath = strings.Trim(filePath, "/")

	err := serveFile(w, filePath, true)
	CheckError(err)
}

// serveIndexPage is handler for GET /
func (h *handler) serveIndexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// make sure session still valid
	err := h.validateSession(r)
	if err != nil {
		newPath := path.Join(h.RootPath, "/login")
		redirectURL := createRedirectURL(newPath, r.URL.String())
		redirectPage(w, r, redirectURL)
		return
	}

	if developmentMode {
		//todo h.prepareTemplates() 没有完成
		if err := h.prepareTemplates(); err != nil {
			log.Printf("error during template preparation : %s", err)
		}
	}

	err = h.templates["index"].Execute(w, h.RootPath)
	CheckError(err)
}

// serveLoginPage is handler for GET /login
func (h *handler) serveLoginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Make sure session is not valid
	err := h.validateSession(r)
	if err == nil {
		redirectURL := path.Join(h.RootPath, "/")
		redirectPage(w, r, redirectURL)
		return
	}

	if developmentMode {
		if err := h.prepareTemplates(); err != nil {
			log.Printf("error during template preparation: %s", err)
		}
	}

	err = h.templates["login"].Execute(w, h.RootPath)
	CheckError(err)
}

// serveVueDemoPage is handler for GET /login
func (h *handler) serveVueDemoPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 啊，你原来是在这里设置的rootPath ，我哭了
	err := h.templates["vue-demo"].Execute(w, h.RootPath)
	CheckError(err)
}
