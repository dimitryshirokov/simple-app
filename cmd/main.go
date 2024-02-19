package main

import (
	"context"
	"github.com/dimitryshirokov/simple-app/internal/config"
	"github.com/dimitryshirokov/simple-app/internal/database"
	"github.com/dimitryshirokov/simple-app/internal/exiter"
	"github.com/dimitryshirokov/simple-app/internal/logger"
	"github.com/dimitryshirokov/simple-app/internal/server"
	"github.com/dimitryshirokov/simple-app/internal/server/handler"
	"github.com/dimitryshirokov/simple-app/internal/service"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := config.NewConfig()
	if err != nil {
		logger.LogError("can't create config", nil, err)
		os.Exit(1)
	}
	dbPool, err := database.CreatePool(ctx, conf)
	if err != nil {
		logger.LogError("can't create database pool", nil, err)
		os.Exit(1)
	}
	calculatorService := service.NewCalculatorService(ctx, conf, dbPool)
	serv := server.NewServer(ctx, conf, map[string]handler.Handler{
		"/addition":    handler.NewAdditionHandler(calculatorService),
		"/subtraction": handler.NewSubtractionHandler(calculatorService),
		"/results":     handler.NewResultsHandler(calculatorService),
	})
	go exiter.HandleExit(ctx, serv)
	serverWait := serv.Run()
	serverWait.Wait()
}
