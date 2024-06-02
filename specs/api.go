package specs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/utils"
)

func getSpecService() *SpecService {
	conf := config.GetConfig()
	service := conf.SpecService.(*SpecService)
	return service
}

func GetAllSpecsHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getSpecService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting spec service", nil)
		return
	}

	specs, err := service.GetAllSpecs()
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", specs)
}

func GetSpecHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getSpecService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting spec service", nil)
		return
	}

	vars := mux.Vars(r)
	specId := vars["specId"]
	exists, err := service.DoesSpecExist(specId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("spec does not exist: %s", specId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	spec, err := service.GetSpec(specId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "success", spec)
}

func CreateSpecHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getSpecService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting spec service", nil)
		return
	}

	var spec TableSpec
	err := json.NewDecoder(r.Body).Decode(&spec)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = service.CreateSpec(&spec)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusCreated, "spec created", spec)
}

func UpdateSpecHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getSpecService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting spec service", nil)
		return
	}

	vars := mux.Vars(r)
	specId := vars["specId"]
	exists, err := service.DoesSpecExist(specId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !exists {
		msg := fmt.Sprintf("spec does not exist: %s", specId)
		utils.CreateApiResponse(w, http.StatusNotFound, msg, nil)
		return
	}

	var spec TableSpec
	err = json.NewDecoder(r.Body).Decode(&spec)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = service.UpdateSpec(specId, &spec)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "spec updated", spec)
}

func DeleteSpecHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	service := getSpecService()
	if service == nil {
		utils.CreateApiResponse(w, http.StatusInternalServerError, "error getting spec service", nil)
		return
	}

	vars := mux.Vars(r)
	specId := vars["specId"]
	err := service.DeleteSpec(specId)
	if err != nil {
		utils.CreateApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.CreateApiResponse(w, http.StatusOK, "spec deleted", nil)
}

func RegisterRoutes() error {
	conf := config.GetConfig()
	router := conf.Router
	router.HandleFunc("/api/v1/specs", GetAllSpecsHandler).Methods("GET")
	router.HandleFunc("/api/v1/specs/{specId}", GetSpecHandler).Methods("GET")
	router.HandleFunc("/api/v1/specs", CreateSpecHandler).Methods("POST")
	router.HandleFunc("/api/v1/specs/{specId}", UpdateSpecHandler).Methods("PUT")
	router.HandleFunc("/api/v1/specs/{specId}", DeleteSpecHandler).Methods("DELETE")
	return nil
}
