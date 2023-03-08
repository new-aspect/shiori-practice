package webserver

import "github.com/new-aspect/shiori-practice/internal/database"

// Config is parameter that used for starting web server
type Config struct {
	DB            database.DB
	DataDir       string
	ServerAddress string
	ServerPort    int
	RootPath      string
	Log           bool
}

// ServeApp serves web interface in specified port
func ServeApp(conf Config) error {
	return nil
}
