package datasets

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

func getService() (common.DatasetService, error) {
	conf := common.GetConfig()
	service := conf.DataSetService.(*DatasetService)
	if service == nil {
		return nil, fmt.Errorf("error getting dataset service")
	}
	return service, nil
}

func GetAllDataSetsHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	ids, err := service.ListDatasetIds()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if ids == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting dataset ids", nil)
		return
	}

	datasets := make([]*common.Dataset, 0)
	for _, id := range ids {
		dataset, err := service.GetDataset(id)
		if err != nil {
			utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		datasets = append(datasets, dataset)
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", datasets)
}

func GetDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	dataSetId := vars["dataSetId"]

	exists, err := service.DoesDatasetExist(dataSetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := "dataset does not exist"
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	dataset, err := service.GetDataset(dataSetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreateApiResponse(w, http.StatusOK, "success", dataset)
}

func CreateDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var dataset common.Dataset
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&dataset)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = service.CreateDataset(&dataset)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreateApiResponse(w, http.StatusCreated, "success", dataset)
}

func DeleteDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	vars := mux.Vars(r)
	dataSetId := vars["dataSetId"]

	err = service.DeleteDataset(dataSetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreateApiResponse(w, http.StatusOK, "success", nil)
}

func AddBatchToDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	vars := mux.Vars(r)
	dataSetId := vars["dataSetId"]
	exists, err := service.DoesDatasetExist(dataSetId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if !exists {
		msg := "dataset does not exist"
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}
	var batch common.Batch
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&batch)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	err = service.AddBatchToDataset(dataSetId, &batch)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.CreateApiResponse(w, http.StatusCreated, "success", nil)
}

func RemoveBatchFromDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	vars := mux.Vars(r)
	dataSetId := vars["dataSetId"]
	batchId := vars["batchId"]
	err = service.RemoveBatchFromDataset(dataSetId, batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", nil)
}

func GetBatchFromDataSetHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service, err := getService()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	vars := mux.Vars(r)
	dataSetId := vars["dataSetId"]
	batchId := vars["batchId"]
	batch, err := service.GetBatchFromDataset(dataSetId, batchId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", batch)
}

func RegisterRoutes() error {
	conf := common.GetConfig()
	router := conf.Router

	router.HandleFunc("/api/v1/datasets", GetAllDataSetsHandler).Methods("GET")
	router.HandleFunc("/api/v1/datasets/{dataSetId}", GetDataSetHandler).Methods("GET")
	router.HandleFunc("/api/v1/datasets", CreateDataSetHandler).Methods("POST")
	router.HandleFunc("/api/v1/datasets/{dataSetId}", DeleteDataSetHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/datasets/{dataSetId}/batches", AddBatchToDataSetHandler).Methods("POST")
	router.HandleFunc("/api/v1/datasets/{dataSetId}/batches/{batchId}", RemoveBatchFromDataSetHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/datasets/{dataSetId}/batches/{batchId}", GetBatchFromDataSetHandler).Methods("GET")
	return nil
}
