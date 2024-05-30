package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/specs"
	"github.com/itaborai83/dsr/utils"
)

var (
	logger *log.Logger
)

func registerServices() error {
	// Specs
	err := specs.RegisterServices()
	if err != nil {
		return fmt.Errorf("error registering spec services: %v", err)
	}
	return nil
}

func registerRoutes() error {
	// Specs
	err := specs.RegisterRoutes()
	if err != nil {
		return fmt.Errorf("error registering spec handlers: %v", err)
	}
	return nil
}

func main() {
	conf := config.GetConfig()
	logger = log.New(os.Stdout, "", log.LstdFlags)

	if !utils.DirExists(conf.DataDir) {
		logger.Fatalf("data directory does not exist: %s\n", conf.DataDir)
		os.Exit(1)
	}

	err := registerServices()
	if err != nil {
		logger.Fatalf("error registering services: %v\n", err)
		os.Exit(1)
	}

	err = registerRoutes()
	if err != nil {
		logger.Fatalf("error registering routes: %v\n", err)
		os.Exit(1)
	}

	logger.Printf("Starting server on %s:%s\n", conf.Host, conf.Port)
	bindPoint := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	err = http.ListenAndServe(bindPoint, conf.Router)
	if err != nil {
		logger.Fatalf("error starting server: %v\n", err)
		os.Exit(1)
	}
}
