package webserver

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/new-aspect/shiori-practice/internal/database"
	cch "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"time"
)

// Config is parameter that used for starting web server
type Config struct {
	DB            database.DB
	DataDir       string
	ServerAddress string
	ServerPort    int
	RootPath      string
	Log           bool
}

// ErrorResponse defines a single HTTP error response
type ErrorResponse struct {
	Code        int
	Body        string
	contentType string
	errorText   string
	Log         bool
}

// responseData will hold response details that we are interested in for logging
type responseData struct {
	status int
	size   int
}

// Wrapper around http.ResponseWriter to be able to catch calls to Write*()
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Logger Log through logrus, 200 will log as info, anything else as an error.
func Logger(r *http.Request, statusCode int, size int) {
	if statusCode == http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"proto":  r.Proto,
			"remote": r.RemoteAddr,
			"reqlen": r.ContentLength,
			"size":   size,
			"status": statusCode,
		}).Info(r.Method, " ", r.RequestURI)
	} else {
		logrus.WithFields(logrus.Fields{
			"proto":  r.Proto,
			"remote": r.RemoteAddr,
			"reqlen": r.ContentLength,
			"size":   size,
			"status": statusCode,
		}).Warn(r.Method, " ", r.RequestURI)
	}
}

func (e *ErrorResponse) Error() string {
	return e.errorText
}

// 感觉这里相当于中间件给http添加一些参数
func (e *ErrorResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.contentType != "" {
		w.Header().Set("Content-Type", e.contentType)
	}
	// todo 这里还没有写完
}

// ServeApp serves web interface in specified port
func ServeApp(cfg Config) error {
	// Create handler
	hdl := handler{
		DB:      cfg.DB,
		DataDir: cfg.DataDir,
		// defaultExpiration 默认过期时间
		// cleanupInterval 清理间隔
		UserCache:    cch.New(time.Hour, 10*time.Minute),
		SessionCache: cch.New(time.Hour, 10*time.Minute),
		ArchiveCache: cch.New(time.Minute, 5*time.Minute),
		RootPath:     cfg.RootPath,
		Log:          cfg.Log,
	}

	hdl.prepareSessionCache()
	hdl.prepareArchiveCache()

	err := hdl.prepareTemplates()
	if err != nil {
		return fmt.Errorf("failed to prepare templates: %v", err)
	}

	// Prepare errors
	var (
		ErrorNotAllow = &ErrorResponse{
			http.StatusMethodNotAllowed,
			"Method is not allowed",
			"text/plain; charset=UTF-8",
			"MethodNotAllowedError",
			cfg.Log,
		}

		ErrorNotFound = &ErrorResponse{
			http.StatusNotFound,
			"Resource Not Found",
			"text/plain; charset=UTF-8",
			"NotFoundError",
			cfg.Log,
		}
	)

	// Create router and register error handlers
	router := httprouter.New()
	router.NotFound = ErrorNotFound
	router.MethodNotAllowed = ErrorNotAllow

	// withLogging will inject our own (compatible) http.ResponseWriter in order
	// to collect details about the answer, i.e. the status code and the size of
	// data in the response. Once done, these are passed further for logging, if
	// relevant.
	withLogging := func(req func(http.ResponseWriter, *http.Request, httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			d := &responseData{
				status: 0,
				size:   0,
			}
			lrw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   d,
			}
			req(&lrw, r, ps)
			if hdl.Log {
				Logger(r, d.status, d.size)
			}
		}
	}

	// jp here means "join path", as in "join route with root path"
	jp := func(route string) string {
		return path.Join(cfg.RootPath, route)
	}

	router.GET(jp("/js/*filepath"), withLogging(hdl.serveJsFile))
	// todo 这里还有很多接口

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, arg interface{}) {
		d := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   d,
		}
		http.Error(&lrw, fmt.Sprint(arg), http.StatusInternalServerError)
		if hdl.Log {
			Logger(r, d.status, d.size)
		}
	}

	// Create server
	url := fmt.Sprintf("%s:%d", cfg.ServerAddress, cfg.ServerPort)
	svr := &http.Server{
		Addr:         url,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: time.Minute,
	}

	//Serve app
	logrus.Infoln("Serve shiori in", url, cfg.RootPath)
	return svr.ListenAndServe()
}
