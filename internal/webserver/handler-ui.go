package webserver

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// serveJsFile is handler for GET /js/*filepath
func (h *handler) serveJsFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// todo 这里没有完成
	fmt.Println("进入到这里了")
}
