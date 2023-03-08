package cmd

import (
	"github.com/new-aspect/shiori-practice/internal/webserver"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

func serveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve web interface for managing bookmarks",
		Long: "Run a simple and performant web server which " +
			"serves the site for managing bookmarks. If --port " +
			"flag is not used, it will use port 8080 by default.",
		Run: serveHandler,
	}

	cmd.Flags().IntP("port", "p", 8080, "Port used by the server")
	cmd.Flags().StringP("address", "a", "", "Address the server listens to")
	cmd.Flags().StringP("webroot", "r", "/", "Root path that used by server")
	cmd.Flags().Bool("log", true, "Print out a non-standard access log")

	return cmd
}

func serveHandler(cmd *cobra.Command, args []string) {
	// Get flags value
	port, _ := cmd.Flags().GetInt("port")
	address, _ := cmd.Flags().GetString("address")
	rootPath, _ := cmd.Flags().GetString("webroot")
	log, _ := cmd.Flags().GetBool("log")

	// Validate root path
	if rootPath == "" {
		rootPath = "/"
	}

	if !strings.HasPrefix(rootPath, "/") {
		rootPath = "/" + rootPath
	}

	if !strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
	}

	// Start server
	serverConfig := webserver.Config{
		DB:            db,
		DataDir:       dataDir,
		ServerAddress: address,
		ServerPort:    port,
		RootPath:      rootPath,
		Log:           log,
	}

	err := webserver.ServeApp(serverConfig)
	if err != nil {
		logrus.Fatalf("Server error %v\n", err)
	}
}
