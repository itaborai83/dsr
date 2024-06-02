package batches

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/utils"
)

func getBatchService() *BatchService {
	conf := config.GetConfig()
	service := conf.BatchService.(*BatchService)
	return service
}

func GetAllBatchIds(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getBatchService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting batch service", nil)
		return
	}

	batchIds, err := service.ListBatchIds()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", batchIds)
}

func GetBatchHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getBatchService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting batch service", nil)
		return
	}

	vars := mux.Vars(r)
	batchId := vars["batchId"]
	exists, err := service.DoesBatchExist(batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("batch does not exist: %s", batchId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	spec, err := service.GetBatch(batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", spec)
}

func CreateBatchHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getBatchService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting batch service", nil)
		return
	}

	// Decode the request body as json
	decoder := json.NewDecoder(r.Body)
	var batch Batch
	err := decoder.Decode(&batch)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, "error decoding request body", nil)
		return
	}

	err = service.CreateBatch(&batch)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusCreated, "success", nil)
}

func RegisterRoutes() error {
	conf := config.GetConfig()
	router := conf.Router
	router.HandleFunc("/api/v1/batches", GetAllBatchIds).Methods("GET")
	router.HandleFunc("/api/v1/batches/{batchId}", GetBatchHandler).Methods("GET")
	return nil
}
