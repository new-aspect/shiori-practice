package cmd

import (
	"context"
	"github.com/new-aspect/shiori-practice/internal/database"
	"github.com/new-aspect/shiori-practice/internal/model"
	"github.com/spf13/cobra"
	"os"
	fp "path/filepath"
)

var (
	db              database.DB
	dataDir         string
	developmentMode bool
)

// ShioriCmd 返回shiori的root cmd
func ShioriCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "shiori",
		Short: "使用go构建的简单书签管理",
	}

	// PersistentPreRun: children of this command will inherit and execute.
	rootCmd.PersistentPreRun = preRunRootHandler
	rootCmd.PersistentFlags().Bool("portable", false, "run shiori in portable mode")
	rootCmd.AddCommand(
		serveCmd(),
	)

	return rootCmd
}

// 初始化数据库
func preRunRootHandler(cmd *cobra.Command, args []string) {
	// Read flag
	var err error
	portableModel, _ := cmd.Flags().GetBool("portable")

	// Get and create data dir
	dataDir, err = getDateDir(portableModel)
	if err != nil {
		_, _ = cError.Printf("Failed to get data dir :%v\n", err)
		os.Exit(1)
	}

	err = os.MkdirAll(dataDir, model.DataDirPerm)
	if err != nil {
		_, _ = cError.Printf("Failed to get data dir :%v\n", err)
		os.Exit(1)
	}

	// Open database
	db, err = openDatabase(cmd.Context())
	if err != nil {
		_, _ = cError.Printf("Failed to open database :%v\n", err)
		os.Exit(1)
	}

	// Migrate 这个表示初始化数据库
	err = db.Migrate()
	if err != nil {
		_, _ = cError.Printf("Error running migration :%s\n", err)
		os.Exit(1)
	}
}

func getDateDir(portableModel bool) (string, error) {
	// If in portable mode, uses directory of executable
	if portableModel {
		exePath, err := os.Executable()
		if err != nil {
			return "", err
		}

		exeDir := fp.Dir(exePath)
		return fp.Join(exeDir, "shiori-data"), nil
	}

	if developmentMode {
		return "dev-data", nil
	}
	return "", nil
}

func openDatabase(ctx context.Context) (database.DB, error) {
	switch dbms, _ := os.LookupEnv("SHIORI_DBMS"); dbms {
	case "mysql":
		return openMysqlDatabase(ctx)
	case "postgresql":
		return openPostgreSQLDatabase(ctx)
	default:
		return openSQLiteDatabase(ctx)
	}
}

func openSQLiteDatabase(ctx context.Context) (database.DB, error) {
	dbPath := fp.Join(dataDir, "shiori.db")
	return database.OpenSQLiteDatabase(ctx, dbPath)
}

func openMysqlDatabase(ctx context.Context) (database.DB, error) {
	return nil, nil
}

func openPostgreSQLDatabase(ctx context.Context) (database.DB, error) {
	return nil, nil
}
