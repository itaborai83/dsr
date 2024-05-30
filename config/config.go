package config

import (
	"os"

	"github.com/gorilla/mux"
)

const (
	DEFAULT_HOST                 = "localhost"
	DEFAULT_PORT                 = "8080"
	DEFAULT_DATA_DIR             = "data"
	DEFAULT_SPECS_DIR            = "specs"
	DEFAULT_BATCHES_DIR          = "batches"
	DEFAULT_DATASETS_DIR         = "datasets"
	DEFAULT_RECONCILIATIONS_DIR  = "reconciliations"
	DEFAULT_SPEC_ENTRY           = "spec.json"
	DEFAULT_BATCH_ENTRY          = "batch.json"
	DEFAULT_DATASET_ENTRY        = "dataset.json"
	DEFAULT_RECONCILIATION_ENTRY = "reconciliation.json"
	ENV_HOST                     = "HOST"
	ENV_PORT                     = "PORT"
	ENV_DATA_DIR                 = "DATA_DIR"
	ENV_SPECS_DIR                = "SPECS_DIR"
	ENV_BATCHES_DIR              = "BATCHES_DIR"
	ENV_DATASETS_DIR             = "DATASETS_DIR"
	ENV_RECONCILIATIONS_DIR      = "RECONCILIATIONS_DIR"
	ENV_SPEC_ENTRY               = "SPEC_ENTRY"
	ENV_BATCH_ENTRY              = "BATCH_ENTRY"
	ENV_DATASET_ENTRY            = "DATASET_ENTRY"
	ENV_RECONCILIATION_ENTRY     = "RECONCILIATION_ENTRY"
)

type Config struct {
	Host                string
	Port                string
	Router              *mux.Router
	DataDir             string
	SpecsDir            string
	BatchesDir          string
	DatasetsDir         string
	ReconciliationsDir  string
	SpecEntry           string
	BatchEntry          string
	DatasetEntry        string
	ReconciliationEntry string
	SpecRepo            interface{}
	SpecService         interface{}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func newConfig() *Config {
	return &Config{
		Host:                getEnv(ENV_HOST, DEFAULT_HOST),
		Port:                getEnv(ENV_PORT, DEFAULT_PORT),
		Router:              mux.NewRouter(),
		DataDir:             getEnv(ENV_DATA_DIR, DEFAULT_DATA_DIR),
		SpecsDir:            getEnv(ENV_SPECS_DIR, DEFAULT_SPECS_DIR),
		BatchesDir:          getEnv(ENV_BATCHES_DIR, DEFAULT_BATCHES_DIR),
		DatasetsDir:         getEnv(ENV_DATASETS_DIR, DEFAULT_DATASETS_DIR),
		ReconciliationsDir:  getEnv(ENV_RECONCILIATIONS_DIR, DEFAULT_RECONCILIATIONS_DIR),
		SpecEntry:           getEnv(ENV_SPEC_ENTRY, DEFAULT_SPEC_ENTRY),
		BatchEntry:          getEnv(ENV_BATCH_ENTRY, DEFAULT_BATCH_ENTRY),
		DatasetEntry:        getEnv(ENV_DATASET_ENTRY, DEFAULT_DATASET_ENTRY),
		ReconciliationEntry: getEnv(ENV_RECONCILIATION_ENTRY, DEFAULT_RECONCILIATION_ENTRY),
		SpecRepo:            nil,
		SpecService:         nil,
	}
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = newConfig()
	}
	return config
}
