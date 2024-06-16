package batches

import (
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/dirrepo"
	"github.com/itaborai83/dsr/hasher"
	"github.com/itaborai83/dsr/utils"
)

type BatchService struct {
	repo        common.BatchRepo
	validator   common.BatchValidator
	specService common.SpecService
}

func RegisterServices() error {
	conf := common.GetConfig()
	dr, err := dirrepo.NewDirRepo(conf.DataDir, conf.DataSetsDir, conf.DataSetEntry, nil)
	if err != nil {
		return fmt.Errorf("error creating dataset repo: %v", err)
	}
	br, err := dirrepo.NewDirRepo(conf.DataDir, conf.BatchesDir, conf.BatchEntry, dr)
	if err != nil {
		return fmt.Errorf("error creating batch repo: %v", err)
	}
	repo, err := NewBatchRepo(br)
	if err != nil {
		return fmt.Errorf("error creating batch repo: %v", err)
	}
	validator := NewBatchValidator(repo, conf.SpecRepo)
	service := NewBatchService(repo, validator, conf.SpecService)
	conf.BatchRepo = repo
	conf.BatchValidator = validator
	conf.BatchService = service
	return nil
}

func NewBatchService(repo common.BatchRepo, validator common.BatchValidator, specService common.SpecService) *BatchService {
	return &BatchService{repo: repo, validator: validator, specService: specService}
}

func (b *BatchService) DoesBatchExist(datasetId, batchId string) (bool, error) {
	exists, err := b.repo.DoesBatchExist(datasetId, batchId)
	return exists, err
}

func (b *BatchService) GetBatch(datasetId, batchId string) (*common.Batch, error) {
	exists, err := b.repo.DoesBatchExist(datasetId, batchId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("batch does not exist: %s", batchId)
	}
	batch, err := b.repo.GetBatch(datasetId, batchId)
	if err != nil {
		return nil, fmt.Errorf("error getting batch: %v", err)
	}
	return batch, nil
}

func (b *BatchService) computeHashes(batch *common.Batch) error {
	h := hasher.NewHasher()
	hashes := make([]uint64, 0)

	// fetch spec
	spec, err := b.specService.GetSpec(batch.SpecId)
	if err != nil {
		return fmt.Errorf("error getting spec: %v", err)
	}

	if spec.KeyColumns == nil || len(spec.KeyColumns) == 0 {
		return fmt.Errorf("no key columns in spec")
	}

	for i := 0; i < batch.RecordCount; i++ {
		// reset hasher
		h.Reset()
		// hash each key field
		for _, keyColumn := range spec.KeyColumns {
			column, ok := batch.Data[keyColumn].([]interface{})
			if !ok {
				return fmt.Errorf("key column not found in batch data: %s", keyColumn)
			}
			if i >= len(column) {
				return fmt.Errorf("key column length mismatch: %s", keyColumn)
			}
			h.Update(column[i])
		}
		hash, err := h.GetHash()
		if err != nil {
			return fmt.Errorf("error getting hash: %v", err)
		}
		hashes = append(hashes, hash)
	}
	batch.RecordHashes = hashes
	batch.RecordCount = len(hashes)
	return nil
}

func (b *BatchService) CreateBatch(datasetId string, batch *common.Batch) error {
	err := b.validator.ValidateCreation(datasetId, batch)
	if err != nil {
		return fmt.Errorf("error validating batch: %v", err)
	}

	err = b.computeHashes(batch)
	if err != nil {
		return fmt.Errorf("error computing hashes: %v", err)
	}

	err = b.repo.CreateBatch(datasetId, batch.Id, batch)
	if err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}
	return nil
}

func (b *BatchService) DeleteBatch(datasetId, batchId string) error {
	err := b.validator.ValidateDeletion(datasetId, batchId)
	if err != nil {
		return fmt.Errorf("error validating deletion: %v", err)
	}
	err = b.repo.DeleteBatch(datasetId, batchId)
	if err != nil {
		return fmt.Errorf("error deleting batch: %v", err)
	}
	return nil
}

func (b *BatchService) ListBatchIds(datasetId string) ([]string, error) {
	ids, err := b.repo.ListBatchIds(datasetId)
	if err != nil {
		return nil, fmt.Errorf("error listing batch ids: %v", err)
	}
	utils.SortStringSlice(ids)
	return ids, nil
}
