package datasets

import (
	"encoding/json"
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

type DatasetRepo struct {
	repo common.Repo
}

func NewDatasetRepo(repo common.Repo) (*DatasetRepo, error) {
	return &DatasetRepo{repo: repo}, nil
}

func (d *DatasetRepo) DoesDatasetExist(id string) (bool, error) {
	return d.repo.DoesEntryExist([]string{id})
}

func (d *DatasetRepo) GetDataset(id string) (*common.Dataset, error) {
	data, err := d.repo.GetEntry([]string{id})
	if err != nil {
		return nil, err
	}
	dataset := &common.Dataset{}
	err = json.Unmarshal(data, dataset)
	if err != nil {
		return nil, err
	}
	return dataset, nil
}

func (d *DatasetRepo) CreateDataset(id string, dataset *common.Dataset) error {
	now := utils.GetNow()
	dataset.CreatedAt = now
	dataset.UpdatedAt = now

	data, err := json.Marshal(dataset)
	if err != nil {
		return err
	}
	err = d.repo.CreateEntry([]string{id}, data)
	if err != nil {
		return fmt.Errorf("error creating dataset: %v", err)
	}
	return nil
}

func (d *DatasetRepo) UpdateDataset(id string, dataset *common.Dataset) error {
	now := utils.GetNow()
	dataset.UpdatedAt = now
	data, err := json.Marshal(dataset)
	if err != nil {
		return err
	}
	err = d.repo.UpdateEntry([]string{id}, data)
	if err != nil {
		return fmt.Errorf("error updating dataset: %v", err)
	}
	return nil
}

func (d *DatasetRepo) DeleteDataset(id string) error {
	err := d.repo.DeleteEntry([]string{id})
	if err != nil {
		return fmt.Errorf("error deleting dataset: %v", err)
	}
	return nil
}

func (d *DatasetRepo) ListDatasetIds() ([]string, error) {
	idStacks, err := d.repo.ListEntryIds(nil)
	if err != nil {
		return nil, fmt.Errorf("error listing dataset ids: %v", err)
	}
	ids := make([]string, len(idStacks))
	for i, idStack := range idStacks {
		lastIdx := len(idStack) - 1
		ids[i] = idStack[lastIdx]
	}
	return ids, nil
}
