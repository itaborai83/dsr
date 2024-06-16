package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mutex  = &sync.Mutex{}
	logger *log.Logger
)

func GetLogger() *log.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}
	return logger
}

func GetNow() string {
	return time.Now().Format(time.RFC3339)
}

func FileExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if stat.IsDir() {
		return false
	}
	return true
}

func DirExists(folderPath string) bool {
	stat, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	if stat.IsDir() {
		return true
	}
	return false
}

func ReadFile(path, fileName string) ([]byte, error) {
	filePath := path + "/" + fileName
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}
	return data, nil
}

func WriteFile(path, fileName string, data []byte) error {
	exists, err := ValidateFolderExists(path)
	if err != nil {
		return fmt.Errorf("error validating folder: '%s'", err)
	}
	if !exists {
		return fmt.Errorf("folder does not exist: '%s'", path)
	}

	filePath := path + "/" + fileName
	// open file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: '%s'", err)
	}
	defer file.Close()
	// write data to file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing file: '%s'", err)
	}
	return nil
}

func DeleteFolder(folderPath string) error {
	err := ValidateId(folderPath)
	if err != nil {
		return fmt.Errorf("error validating folder path: %s", err)
	}
	exists, err := ValidateFolderExists(folderPath)
	if err != nil {
		return fmt.Errorf("error validating folder: %s", err)
	}
	if !exists {
		return fmt.Errorf("folder does not exist")
	}
	err = os.RemoveAll(folderPath)
	if err != nil {
		return fmt.Errorf("error deleting folder: %s", err)
	}
	return nil
}

func ListFolders(basePath string) ([]string, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err)
	}

	subFolders := make([]string, 0)
	for _, file := range files {
		// is it a folder
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			subFolders = append(subFolders, file.Name())
		}
	}
	return subFolders, nil
}

func ValidateFileExists(basePath string, fileName string) (bool, error) {
	err := ValidateId(fileName)
	if err != nil {
		return false, fmt.Errorf("error validating file name: %s", err)
	}

	exists, err := ValidateFolderExists(basePath)
	if err != nil {
		return false, fmt.Errorf("error validating folder: %s", err)
	}
	if !exists {
		return false, nil
	}

	filePath := basePath + "/" + fileName
	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error checking file: %s", err)
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	if fileInfo.IsDir() {
		return false, fmt.Errorf("file is a folder")
	}

	return exists, nil
}

func ValidateFolderExists(folderPath string) (bool, error) {
	err := ValidateId(folderPath)
	if err != nil {
		return false, err
	}
	fileInfo, err := os.Stat(folderPath)
	os.IsNotExist(err)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error checking folder: %s", err)
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	if !fileInfo.IsDir() {
		return false, fmt.Errorf("folder is a file")
	}
	return true, nil
}

func EnsureFolderExists(basePath string, folderName string) error {
	err := ValidateId(folderName)
	if err != nil {
		return fmt.Errorf("invalid folder name: '%s'", folderName)
	}

	exists, err := ValidateFolderExists(basePath)
	if err != nil {
		return fmt.Errorf("error validating folder: %s", err)
	}
	if !exists {
		return fmt.Errorf("base folder does not exist")
	}

	folderPath := basePath + "/" + folderName
	exists, err = ValidateFolderExists(folderPath)
	if err != nil {
		return fmt.Errorf("error validating folder: %s", err)
	}

	if exists {
		return nil
	}

	err = os.Mkdir(folderPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating folder: %s", err)
	}
	return nil
}

func HasPathTraversal(fileName string) bool {
	if strings.Contains(fileName, "..") {
		return true
	}

	if strings.HasPrefix(fileName, "/") || strings.HasPrefix(fileName, "\\") {
		return true
	}
	return false
}

func ValidateId(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if strings.Contains(id, " ") {
		return fmt.Errorf("id contains spaces")
	}
	if HasPathTraversal(id) {
		return fmt.Errorf("id contains path traversal")
	}
	return nil
}

func ReverseStringList(list []string) []string {
	result := make([]string, len(list))
	for i := 0; i < len(list); i++ {
		result[i] = list[len(list)-i-1]
	}
	return result
}

func CopyPushString(list []string, item string) []string {
	result := make([]string, len(list)+1)
	copy(result, list)
	result[len(list)] = item
	return result
}

func LogRequest(r *http.Request) {
	logger = GetLogger()
	logger.Printf("Request: %s %s\n", r.Method, r.URL)
	logger.Printf("Headers: %v\n", r.Header)
}

func CreateApiResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	apiResponse := struct {
		StatusCode int         `json:"status"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
	}{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(apiResponse)
}

func SortStringSlice(list []string) error {
	if list == nil {
		return fmt.Errorf("list is nil")
	}
	// sort strings in place using selection sort
	for i := 0; i < len(list); i++ {
		// find the minimum element in the unsorted part of the list
		minIndex := i
		for j := i + 1; j < len(list); j++ {
			if list[j] < list[minIndex] {
				minIndex = j
			}
		}
		// swap the minimum element with the first element in the unsorted part
		list[i], list[minIndex] = list[minIndex], list[i]
	}
	return nil
}

func RemoveStringFromSlice(list []string, item string) []string {
	result := make([]string, 0)
	for _, value := range list {
		if value != item {
			result = append(result, value)
		}
	}
	return result
}
