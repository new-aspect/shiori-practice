package main

import (
	// Database driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/new-aspect/shiori-practice/internal/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cmd.ShioriCmd().Execute()
	if err != nil {
		logrus.Fatalln(err)
	}
}
