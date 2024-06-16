package dirrepo

import (
	"fmt"
	"strings"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

type DirRepo struct {
	basePath     string
	parentRepo   *DirRepo
	entityName   string
	entryName    string
	nestingLevel int
}

func NewDirRepo(basePath, entityName, entryName string, parentRepo *DirRepo) (*DirRepo, error) {
	exists, err := utils.ValidateFolderExists(basePath)
	if err != nil {
		return nil, fmt.Errorf("error checking base path '%s': %v", basePath, err)
	}
	if !exists {
		return nil, fmt.Errorf("base path '%s' does not exist", basePath)
	}

	err = utils.ValidateId(entityName)
	if err != nil {
		return nil, fmt.Errorf("invalid entitys name '%s': %v", entityName, err)
	}

	err = utils.ValidateId(entryName)
	if err != nil {
		return nil, fmt.Errorf("invalid entry name '%s': %v", entryName, err)
	}

	if parentRepo != nil {
		if parentRepo.basePath != basePath {
			return nil, fmt.Errorf("base path mismatch: '%s' != '%s'", parentRepo.basePath, basePath)
		}
	}

	nestingLevel := 0
	repo := parentRepo
	for repo != nil {
		repo = repo.parentRepo
		nestingLevel++
	}

	result := &DirRepo{
		basePath:     basePath,
		parentRepo:   parentRepo,
		entityName:   entityName,
		entryName:    entryName,
		nestingLevel: nestingLevel,
	}
	return result, nil
}

func (d *DirRepo) isValidIdStack(idStack []string) error {
	if len(idStack) == 0 {
		return fmt.Errorf("id stack is empty")
	}
	if len(idStack) != d.nestingLevel+1 {
		return fmt.Errorf("invalid id stack")
	}
	for _, id := range idStack {
		err := utils.ValidateId(id)
		if err != nil {
			return fmt.Errorf("invalid id '%s': %v", id, err)
		}
	}
	return nil
}

func (d *DirRepo) getEntityPath(idStack []string) (string, error) {
	parentRepo := d.GetParentRepo().(*DirRepo)
	if parentRepo == nil {
		result := fmt.Sprintf("%s/%s", d.basePath, d.entityName)
		return result, nil
	}

	parentIds := idStack[:d.nestingLevel]
	err := parentRepo.isValidIdStack(parentIds)
	if err != nil {
		return "", fmt.Errorf("invalid id stack: %v", err)
	}

	parentRepo = d.GetParentRepo().(*DirRepo)
	pathElements := make([]string, 1)
	pathElements[0] = d.basePath
	repoChain := parentRepo.getRepoChain()
	for idx, repo := range repoChain {
		pathElements = append(pathElements, repo.entityName)
		pathElements = append(pathElements, idStack[idx])
	}
	pathElements = append(pathElements, d.entityName)
	path := strings.Join(pathElements, "/")
	return path, nil
}

func (d *DirRepo) getEntryPath(idStack []string) (string, string, error) {
	err := d.isValidIdStack(idStack)
	if err != nil {
		return "", "", fmt.Errorf("invalid id stack: %v", err)
	}
	entityPath, err := d.getEntityPath(idStack)
	if err != nil {
		return "", "", fmt.Errorf("error getting entity path: %v", err)
	}
	lastId := idStack[len(idStack)-1]
	entryPath := entityPath + "/" + lastId
	return entryPath, d.entryName, nil
}

func (d *DirRepo) GetParentRepo() common.Repo {
	return d.parentRepo
}

func (d *DirRepo) HasParentRepo() bool {
	return d.parentRepo != nil
}

func (d *DirRepo) DoesEntryExist(id []string) (bool, error) {
	err := d.isValidIdStack(id)
	if err != nil {
		return false, err
	}

	entryPath, entryName, err := d.getEntryPath(id)
	if err != nil {
		return false, fmt.Errorf("error getting entry path: %v", err)
	}

	exists, err := utils.ValidateFileExists(entryPath, entryName)
	if err != nil {
		return false, nil
	}
	return exists, nil
}

func (d *DirRepo) GetEntry(id []string) ([]byte, error) {
	exists, err := d.DoesEntryExist(id)
	if err != nil {
		return nil, fmt.Errorf("error checking entry: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("entry does not exist")
	}
	entryPath, entryName, err := d.getEntryPath(id)
	if err != nil {
		return nil, fmt.Errorf("error getting entry path: %v", err)
	}
	data, err := utils.ReadFile(entryPath, entryName)
	if err != nil {
		return nil, fmt.Errorf("error reading entry: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("entry is empty")
	}
	return data, nil
}

func (d *DirRepo) getRepoChain() []*DirRepo {
	repoChain := make([]*DirRepo, 0)
	repo := d
	for repo != nil {
		repoChain = append(repoChain, repo)
		repo = repo.GetParentRepo().(*DirRepo)
	}
	// reverse the chain in place
	for i := 0; i < len(repoChain)/2; i++ {
		j := len(repoChain) - i - 1
		repoChain[i], repoChain[j] = repoChain[j], repoChain[i]
	}
	return repoChain
}

func (d *DirRepo) ensureEntryPathExists(idStack []string) error {
	parentFolder := d.basePath
	repoChain := d.getRepoChain()
	for idx, repo := range repoChain {
		// ensure entity folder gets created
		err := utils.EnsureFolderExists(parentFolder, repo.entityName)
		if err != nil {
			return fmt.Errorf("error ensuring entity folder exists: %v", err)
		}
		// ensure entry folder gets created
		id := idStack[idx]
		err = utils.EnsureFolderExists(parentFolder+"/"+repo.entityName, id)
		if err != nil {
			return fmt.Errorf("error ensuring entry folder exists: %v", err)
		}
		// set parent folder for next iteration
		parentFolder = parentFolder + "/" + repo.entityName + "/" + id
	}
	return nil
}

func (d *DirRepo) saveEntry(id []string, data []byte) error {
	err := d.isValidIdStack(id)
	if err != nil {
		return fmt.Errorf("invalid id stack: %v", err)
	}

	if data == nil {
		return fmt.Errorf("data is nil")
	}

	err = d.ensureEntryPathExists(id)
	if err != nil {
		return fmt.Errorf("error ensuring entry path exists: %v", err)
	}

	entryPath, entryName, err := d.getEntryPath(id)
	if err != nil {
		return fmt.Errorf("error getting entry path: %v", err)
	}

	err = utils.WriteFile(entryPath, entryName, data)
	if err != nil {
		return fmt.Errorf("error writing entry: %v", err)
	}
	return nil
}

func (d *DirRepo) CreateEntry(id []string, data []byte) error {
	exists, err := d.DoesEntryExist(id)
	if err != nil {
		return fmt.Errorf("error checking entry: %v", err)
	}
	if exists {
		return fmt.Errorf("entry already exists")
	}
	err = d.saveEntry(id, data)
	if err != nil {
		return fmt.Errorf("error saving entry: %v", err)
	}
	return nil
}

func (d *DirRepo) UpdateEntry(id []string, data []byte) error {
	exists, err := d.DoesEntryExist(id)
	if err != nil {
		return fmt.Errorf("error checking entry: %v", err)
	}
	if !exists {
		return fmt.Errorf("entry does not exist")
	}
	err = d.saveEntry(id, data)
	if err != nil {
		return fmt.Errorf("error saving entry: %v", err)
	}
	return nil
}

func (d *DirRepo) DeleteEntry(id []string) error {
	err := d.isValidIdStack(id)
	if err != nil {
		return fmt.Errorf("invalid id stack: %v", err)
	}

	exists, err := d.DoesEntryExist(id)
	if err != nil {
		return fmt.Errorf("error checking entry: %v", err)
	}
	if !exists {
		return fmt.Errorf("entry does not exist")
	}

	entryPath, _, err := d.getEntryPath(id)
	if err != nil {
		return fmt.Errorf("error getting entry path: %v", err)
	}
	utils.DeleteFolder(entryPath)
	return nil
}

func (d *DirRepo) ListEntryIds(parentId []string) ([][]string, error) {
	parentFolder := d.basePath + "/" + d.entityName
	if d.HasParentRepo() {
		repoChain := d.parentRepo.getRepoChain()
		for idx, repo := range repoChain {
			id := parentId[idx]
			parentFolder += "/" + repo.entityName + "/" + id
		}
		parentFolder += "/" + d.entityName
	}
	folders, err := utils.ListFolders(parentFolder)
	if err != nil {
		return nil, fmt.Errorf("error listing entry ids: %v", err)
	}
	result := make([][]string, 0)
	for _, folder := range folders {
		idStack := utils.CopyPushString(parentId, folder)
		result = append(result, idStack)
	}
	return result, nil
}
