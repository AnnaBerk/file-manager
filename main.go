package main

import (
	"f-manager/config"
	"f-manager/controller"
	"f-manager/pkg/httpserver"
	"f-manager/pkg/psql"
	"f-manager/pkg/validator"
	"f-manager/repo/pgdb"
	"f-manager/service"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.LoadConfig()

	storage, err := psql.New(cfg.ConnectionString)
	if err != nil {
		log.Fatalln("failed to init storage", err)
		os.Exit(1)
	}
	defer storage.Close()

	repo := pgdb.NewFileManagerRepo(storage)

	services := service.NewFileManagerService(repo)

	handler := echo.New()

	handler.Validator = validator.InitializeValidator()
	controller.NewRouter(handler, services)

	// HTTP server
	log.Printf("Starting http server...")
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatalf("app - Run - httpServer.Notify: %w", err)
	}

	// Graceful shutdown
	log.Printf("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Fatalf("app - Run - httpServer.Shutdown: %w", err)
	}

}
