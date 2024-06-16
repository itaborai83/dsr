package batches

import (
	"encoding/json"
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

type BatchRepo struct {
	repo common.Repo
}

func NewBatchRepo(repo common.Repo) (*BatchRepo, error) {
	return &BatchRepo{repo: repo}, nil
}

func (b *BatchRepo) DoesBatchExist(datasetId, id string) (bool, error) {
	exists, err := b.repo.DoesEntryExist([]string{datasetId, id})
	if err != nil {
		return false, fmt.Errorf("error checking if batch exists: %v", err)
	}
	return exists, nil
}

func (b *BatchRepo) GetBatch(datasetId, id string) (*common.Batch, error) {
	data, err := b.repo.GetEntry([]string{datasetId, id})
	if err != nil {
		return nil, err
	}
	batch := &common.Batch{}
	err = json.Unmarshal(data, batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (b *BatchRepo) CreateBatch(datasetId, id string, batch *common.Batch) error {
	now := utils.GetNow()
	batch.CreatedAt = now
	batch.UpdatedAt = now

	data, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	err = b.repo.CreateEntry([]string{datasetId, id}, data)
	if err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}
	return nil
}

func (b *BatchRepo) UpdateBatch(datasetId, id string, batch *common.Batch) error {
	now := utils.GetNow()
	batch.UpdatedAt = now
	data, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	err = b.repo.UpdateEntry([]string{datasetId, id}, data)
	if err != nil {
		return fmt.Errorf("error updating batch: %v", err)
	}
	return nil
}

func (b *BatchRepo) DeleteBatch(datasetId, id string) error {
	err := b.repo.DeleteEntry([]string{datasetId, id})
	if err != nil {
		return fmt.Errorf("error deleting batch: %v", err)
	}
	return nil
}

func (b *BatchRepo) ListBatchIds(datasetId string) ([]string, error) {
	idStacks, err := b.repo.ListEntryIds([]string{datasetId})
	if err != nil {
		return nil, fmt.Errorf("error listing batch ids: %v", err)
	}
	ids := make([]string, len(idStacks))
	for i, idStack := range idStacks {
		lastIndex := len(idStack) - 1
		ids[i] = idStack[lastIndex]
	}
	return ids, nil
}
