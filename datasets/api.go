package datasets

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itaborai83/dsr/batches"
	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/utils"
)

func getDataSetsService() *DataSetService {
	conf := config.GetConfig()
	service := conf.DataSetService.(*DataSetService)
	return service
}

func GetAllDataSetsHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting data sets service", nil)
		return
	}

	dataSets, err := service.ListDataSets()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", dataSets)
}

func GetDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting batch service", nil)
		return
	}

	vars := mux.Vars(r)
	batchId := vars["datasetId"]
	exists, err := service.DoesDataSetExist(batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("data set does not exist: %s", batchId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	dataSet, err := service.GetDataSet(batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", dataSet)
}

func CreateDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting data sets service", nil)
		return
	}

	// Decode the request body as json
	decoder := json.NewDecoder(r.Body)
	var dataSet DataSet
	err := decoder.Decode(&dataSet)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, "error decoding request body", nil)
		return
	}

	err = service.CreateDataSet(&dataSet)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusCreated, "success", dataSet)
}

func DeleteDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting data sets service", nil)
		return
	}

	vars := mux.Vars(r)
	batchId := vars["datasetId"]
	exists, err := service.DoesDataSetExist(batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("data set does not exist: %s", batchId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	err = service.DeleteDataSet(batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", nil)
}

func AddBatchToDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting data sets service", nil)
		return
	}

	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	exists, err := service.DoesDataSetExist(datasetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("data set does not exist: %s", datasetId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var batch batches.Batch
	err = decoder.Decode(&batch)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, "error decoding request body", nil)
		return
	}

	err = service.AddBatchToDataSet(datasetId, &batch)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusCreated, "success", batch)
}

func DeleteBatchFromDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting data sets service", nil)
		return
	}

	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	exists, err := service.DoesDataSetExist(datasetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("data set does not exist: %s", datasetId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	batchId := vars["batchId"]
	err = service.RemoveBatchFromDataSet(datasetId, batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", nil)
}

func GetBatchFromDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getDataSetsService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting data sets service", nil)
		return
	}

	vars := mux.Vars(r)
	datasetId := vars["datasetId"]
	exists, err := service.DoesDataSetExist(datasetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("data set does not exist: %s", datasetId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	batchId := vars["batchId"]
	batch, err := service.GetBatchFromDataSet(datasetId, batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", batch)
}

func RegisterRoutes() error {
	conf := config.GetConfig()
	router := conf.Router
	router.HandleFunc("/api/v1/datasets", GetAllDataSetsHandler).Methods("GET")
	router.HandleFunc("/api/v1/datasets/{datasetId}", GetDataSetHandler).Methods("GET")
	router.HandleFunc("/api/v1/datasets", CreateDataSetHandler).Methods("POST")
	router.HandleFunc("/api/v1/datasets/{datasetId}", DeleteDataSetHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/datasets/{datasetId}/batches", AddBatchToDataSetHandler).Methods("POST")
	router.HandleFunc("/api/v1/datasets/{datasetId}/batches/{batchId}", DeleteBatchFromDataSetHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/datasets/{datasetId}/batches/{batchId}", GetBatchFromDataSetHandler).Methods("GET")

	return nil
}
