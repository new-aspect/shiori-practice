package webserver

import (
	"github.com/new-aspect/shiori-practice/internal/database"
	cch "github.com/patrickmn/go-cache"
	"html/template"
)

// Handler is handler for serving the web interface
type handler struct {
	DB           database.DB
	DataDir      string
	RootPath     string
	UserCache    *cch.Cache
	SessionCache *cch.Cache
	ArchiveCache *cch.Cache
	Log          bool
	templates    map[string]*template.Template
}

var developmentMode = false

// todo 这里有三个方法没有实现
func (h *handler) prepareSessionCache() {

}

func (h *handler) prepareArchiveCache() {

}

func (h *handler) prepareTemplates() error {
	return nil
}
