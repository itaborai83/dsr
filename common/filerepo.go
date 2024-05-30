package common

import (
	"fmt"

	"github.com/itaborai83/dsr/utils"
)

type FileRepo struct {
	baseDir       string
	entryFileName string
}

func NewFileRepo(baseDir, entryFileName string) (*FileRepo, error) {
	// check if baseDir exists
	err := utils.ValidateFolderExists(baseDir)
	if err != nil {
		return nil, fmt.Errorf("error validating base directory exists: %v", err)
	}
	// check if entryName is valid

	if utils.HasPathTraversal(entryFileName) {
		return nil, fmt.Errorf("invalid entry name: %v", entryFileName)
	}
	result := &FileRepo{baseDir: baseDir, entryFileName: entryFileName}
	return result, nil
}

func (f *FileRepo) DoesEntryExist(entryId string) (bool, error) {
	// validate entry id
	if utils.HasPathTraversal(entryId) {
		return false, fmt.Errorf("invalid entry id: %s", entryId)
	}
	// validate entry dir ex
	entryDir := f.baseDir + "/" + entryId
	err := utils.ValidateFolderExists(entryDir)
	if err != nil {
		return false, nil
	}
	// validate entry file exists
	entryFile := entryDir + "/" + f.entryFileName
	return utils.FileExists(entryFile), nil
}

func (f *FileRepo) GetEntry(entryId string) ([]byte, error) {
	exists, err := f.DoesEntryExist(entryId)
	if err != nil {
		return nil, fmt.Errorf("error checking entry: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("entry does not exist")
	}
	entryFile := f.baseDir + "/" + entryId + "/" + f.entryFileName
	data, err := utils.ReadFile(entryFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return data, nil
}

func (f *FileRepo) CreateEntry(entryId string, data []byte) error {
	exists, err := f.DoesEntryExist(entryId)
	if err != nil {
		return fmt.Errorf("error checking entry: %v", err)
	}
	if exists {
		return fmt.Errorf("entry already exists")
	}
	entryDir := f.baseDir + "/" + entryId
	err = utils.EnsureFolderExists(f.baseDir, entryId)
	if err != nil {
		return fmt.Errorf("error creating entry directory: %v", err)
	}
	entryFile := entryDir + "/" + f.entryFileName
	err = utils.WriteFile(entryFile, data)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func (f *FileRepo) UpdateEntry(entryId string, data []byte) error {
	exists, err := f.DoesEntryExist(entryId)
	if err != nil {
		return fmt.Errorf("error checking entry: %v", err)
	}
	if !exists {
		return fmt.Errorf("entry does not exist")
	}
	entryFile := f.baseDir + "/" + entryId + "/" + f.entryFileName
	err = utils.WriteFile(entryFile, data)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func (f *FileRepo) DeleteEntry(entryId string) error {
	exists, err := f.DoesEntryExist(entryId)
	if err != nil {
		return fmt.Errorf("error checking entry: %v", err)
	}
	if !exists {
		return fmt.Errorf("entry does not exist")
	}
	entryDir := f.baseDir + "/" + entryId
	err = utils.DeleteFolder(entryDir)
	if err != nil {
		return fmt.Errorf("error deleting entry: %v", err)
	}
	return nil
}

func (f *FileRepo) ListEntryIds() ([]string, error) {
	entries, err := utils.ListFolders(f.baseDir)
	if err != nil {
		return nil, fmt.Errorf("error listing entries: %v", err)
	}
	return entries, nil
}

/*




func (f *FileRepo) ListEntryIds() ([]string, error)
*/
