package batches

import (
	"encoding/json"
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/specs"
	"github.com/itaborai83/dsr/utils"
)

type BatchService struct {
	repo        common.Repo
	specService *specs.SpecService
}

func RegisterServices() error {
	conf := config.GetConfig()
	err := utils.EnsureFolderExists(conf.DataDir, conf.BatchesDir)
	if err != nil {
		return fmt.Errorf("error ensuring batches directory exists: %v", err)
	}
	batchesDir := conf.DataDir + "/" + conf.BatchesDir
	repo, err := common.NewFileRepo(batchesDir, conf.BatchEntry)
	if err != nil {
		return fmt.Errorf("error creating spec repo: %v", err)
	}
	specService := conf.SpecService.(*specs.SpecService)
	if specService == nil {
		return fmt.Errorf("error getting previously registered spec service")
	}
	service := NewBatchService(repo, specService)
	conf.BatchRepo = repo
	conf.BatchService = service
	return nil
}

func NewBatchService(repo common.Repo, specService *specs.SpecService) *BatchService {
	return &BatchService{repo: repo, specService: specService}
}

func (b *BatchService) validateAgainstSpec(batch Batch) error {

	exists, err := b.specService.DoesSpecExist(batch.SpecId)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", batch.SpecId)
	}

	spec, err := b.specService.GetSpec(batch.SpecId)
	if err != nil {
		return fmt.Errorf("error getting spec: %v", err)
	}

	conforms := spec.ConformsTo(batch.Data)
	if !conforms {
		return fmt.Errorf("batch data does not conform to spec: %s", batch.SpecId)
	}
	return nil
}

func (b *BatchService) DoesBatchExist(batchId string) (bool, error) {
	return b.repo.DoesEntryExist(batchId)
}

func (b *BatchService) GetBatch(batchId string) (*Batch, error) {
	data, err := b.repo.GetEntry(batchId)
	if err != nil {
		return nil, fmt.Errorf("error getting batch: %v", err)
	}
	var batch Batch
	err = json.Unmarshal(data, &batch)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling batch: %v", err)
	}
	return &batch, nil
}

func (b *BatchService) CreateBatch(batch *Batch) error {
	err := ValidateBatch(batch)
	if err != nil {
		return fmt.Errorf("error validating batch: %v", err)
	}
	exists, err := b.DoesBatchExist(batch.Id)
	if err != nil {
		return fmt.Errorf("error checking if batch exists: %v", err)
	}
	if exists {
		return fmt.Errorf("batch already exists: %s", batch.Id)
	}

	err = b.validateAgainstSpec(*batch)
	if err != nil {
		return fmt.Errorf("error validating batch against spec: %v", err)
	}

	now := utils.GetNow()
	batch.CreatedAt = now
	batch.UpdatedAt = now
	data, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("error marshalling batch: %v", err)
	}
	err = b.repo.CreateEntry(batch.Id, data)
	if err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}
	return nil
}

func (b *BatchService) UpdateBatch(batchId string, batch *Batch) error {
	if batchId != batch.Id {
		return fmt.Errorf("batch id mismatch")
	}

	err := ValidateBatch(batch)
	if err != nil {
		return fmt.Errorf("error validating batch: %v", err)
	}

	exists, err := b.DoesBatchExist(batchId)
	if err != nil {
		return fmt.Errorf("error checking if batch exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("batch does not exist: %s", batchId)
	}

	err = b.validateAgainstSpec(*batch)
	if err != nil {
		return fmt.Errorf("error validating batch against spec: %v", err)
	}

	now := utils.GetNow()
	batch.UpdatedAt = now
	data, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("error marshalling batch: %v", err)
	}
	err = b.repo.UpdateEntry(batchId, data)
	if err != nil {
		return fmt.Errorf("error updating batch: %v", err)
	}
	return nil
}

func (b *BatchService) DeleteBatch(batchId string) error {
	exists, err := b.DoesBatchExist(batchId)
	if err != nil {
		return fmt.Errorf("error checking if batch exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("batch does not exist: %s", batchId)
	}
	err = b.repo.DeleteEntry(batchId)
	if err != nil {
		return fmt.Errorf("error deleting batch: %v", err)
	}
	return nil
}

func (b *BatchService) ListBatchIds() ([]string, error) {
	ids, err := b.repo.ListEntryIds()
	if err != nil {
		return nil, fmt.Errorf("error listing batch ids: %v", err)
	}
	return ids, nil
}
