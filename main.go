package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/itaborai83/dsr/batches"
	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/datasets"
	"github.com/itaborai83/dsr/specs"
	"github.com/itaborai83/dsr/utils"
)

func registerServices() error {
	// Specs
	err := specs.RegisterServices()
	if err != nil {
		return fmt.Errorf("error registering spec services: %v", err)
	}
	// Batches
	err = batches.RegisterServices()
	if err != nil {
		return fmt.Errorf("error registering batch services: %v", err)
	}
	// DataSets
	err = datasets.RegisterServices()
	if err != nil {
		return fmt.Errorf("error registering dataset services: %v", err)
	}
	return nil
}

func registerRoutes() error {
	// healthcheck
	conf := common.GetConfig()
	conf.Router.HandleFunc("/healthcheck", HealthCheck)
	// Specs
	err := specs.RegisterRoutes()
	if err != nil {
		return fmt.Errorf("error registering spec handlers: %v", err)
	}
	// Batches - as batches are a weak entity, the routes are registered in the datasets package
	// DataSets
	err = datasets.RegisterRoutes()
	if err != nil {
		return fmt.Errorf("error registering dataset handlers: %v", err)
	}
	return nil
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	conf := common.GetConfig()
	logger := utils.GetLogger()

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
