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

func ReadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}
	return data, nil
}

func WriteFile(filePath string, data []byte) error {
	// open file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer file.Close()
	// write data to file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}
	return nil
}

func DeleteFolder(folderPath string) error {
	err := os.RemoveAll(folderPath)
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

func ValidateFileExists(basePath string, fileName string) error {
	err := ValidateId(fileName)
	if err != nil {
		return err
	}

	err = ValidateFolderExists(basePath)
	if err != nil {
		return fmt.Errorf("base path folder does not exist: %s", err)
	}

	filePath := fmt.Sprintf("%s/%s", basePath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}
	return nil
}

func ValidateFolderExists(folderPath string) error {
	err := ValidateId(folderPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist")
	}
	return nil
}

func EnsureFolderExists(basePath string, folderName string) error {
	err := ValidateId(folderName)
	if err != nil {
		return fmt.Errorf("error validating folder name: %s", err)
	}

	err = ValidateFolderExists(basePath)
	if err != nil {
		return fmt.Errorf("base path does not exist: %s", err)
	}

	folderPath := fmt.Sprintf("%s/%s", basePath, folderName)
	err = ValidateFolderExists(folderPath)
	if err == nil {
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
