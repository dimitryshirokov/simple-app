package main

import (
	"context"
	"github.com/dimitryshirokov/simple-app/internal/config"
	"github.com/dimitryshirokov/simple-app/internal/database"
	"github.com/dimitryshirokov/simple-app/internal/logger"
	"github.com/dimitryshirokov/simple-app/internal/migrations"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := config.NewConfig()
	if err != nil {
		logger.LogError("can't create config", nil, err)
		os.Exit(1)
	}
	c.DbMaxConnections = 1
	c.DbMinConnections = 1
	conn, err := database.CreatePool(ctx, c)
	if err != nil {
		logger.LogError("can't create database connection", nil, err)
		os.Exit(1)
	}
	defer conn.Close()
	mr := migrations.NewMigratorRunner(conn)
	if len(os.Args) == 2 && os.Args[1] == "generate" {
		err := mr.CreateMigration()
		if err != nil {
			logger.LogError("can't generate migration", nil, err)
			os.Exit(1)
		}
		log.Println("migration created")
		os.Exit(0)
	}
	err = mr.Migrate()
	if err != nil {
		logger.LogError("can't migrate to current version", nil, err)
		os.Exit(1)
	}
}
