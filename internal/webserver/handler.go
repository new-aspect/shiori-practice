package webserver

import (
	"fmt"
	"github.com/new-aspect/shiori-practice/internal/database"
	"github.com/new-aspect/shiori-practice/internal/model"
	cch "github.com/patrickmn/go-cache"
	"html/template"
	"net/http"
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
	// Prepare variable
	var err error
	h.templates = make(map[string]*template.Template)

	// Prepare func map
	funcMap := template.FuncMap{
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// Create template for login, index and content
	for _, name := range []string{"login", "index", "content", "v1", "v2"} {
		h.templates[name], err = createTemplate(name+".html", funcMap)
		if err != nil {
			return err
		}
	}

	// Create template for archive overlay
	h.templates["archive"], err = template.New("archive").Delims("$$", "$$").Parse(
		`<div id="shiori-archive-header">
		<p id="shiori-logo"><span>栞</span>shiori</p>
		<div class="spacer"></div>
		<a href="$$.URL$$" target="_blank">View Original</a>
		$$if .HasContent$$
		<a href="/bookmark/$$.ID$$/content">View Readable</a>
		$$end$$
		</div>`)
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) getSessionID(r *http.Request) string {
	// Try to get session ID from the header
	sessionID := r.Header.Get("X-Session-Id")

	// If not, try it form cookie
	if sessionID == "" {
		cookie, err := r.Cookie("session-id")
		if err != nil {
			return ""
		}

		sessionID = cookie.Value
	}

	return sessionID
}

// validateSession checks whether user session is still valid or not
func (h *handler) validateSession(r *http.Request) error {
	sessionID := h.getSessionID(r)
	if sessionID == "" {
		return fmt.Errorf("session is not exist")
	}

	// Make sure session is not expired yet
	val, found := h.SessionCache.Get(sessionID)
	if !found {
		return fmt.Errorf("session has been expried")
	}

	// If this is not get request, make sure it's over
	if r.Method != "" && r.Method != http.MethodGet {
		if account := val.(model.Account); !account.Owner {
			return fmt.Errorf("account level is not sufficient")
		}
	}

	return nil
}
