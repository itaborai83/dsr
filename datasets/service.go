package datasets

import (
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/dirrepo"
	"github.com/itaborai83/dsr/utils"
)

type DatasetService struct {
	repo         common.DatasetRepo
	validator    common.DatasetValidator
	batchService common.BatchService
}

func RegisterServices() error {
	conf := common.GetConfig()
	dr, err := dirrepo.NewDirRepo(conf.DataDir, conf.DataSetsDir, conf.DataSetEntry, nil)
	if err != nil {
		return fmt.Errorf("error creating dataset repo: %v", err)
	}
	repo, err := NewDatasetRepo(dr)
	if err != nil {
		return fmt.Errorf("error creating dataset repo: %v", err)
	}
	validator := NewDatasetValidator(repo, conf.SpecRepo)
	service := NewDatasetService(repo, validator, conf.BatchService)
	conf.DataSetService = service
	return nil
}

func NewDatasetService(repo common.DatasetRepo, validator common.DatasetValidator, batchService common.BatchService) *DatasetService {
	return &DatasetService{repo: repo, validator: validator, batchService: batchService}
}

func (s *DatasetService) DoesDatasetExist(dataSetId string) (bool, error) {
	exists, err := s.repo.DoesDatasetExist(dataSetId)
	if err != nil {
		return false, fmt.Errorf("error checking if dataset exists: %v", err)
	}
	return exists, nil
}

func (s *DatasetService) GetDataset(dataSetId string) (*common.Dataset, error) {
	dataset, err := s.repo.GetDataset(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error getting dataset: %v", err)
	}
	return dataset, nil
}

func (s *DatasetService) CreateDataset(dataSet *common.Dataset) error {
	err := s.validator.ValidateCreation(dataSet)
	if err != nil {
		return err
	}
	return s.repo.CreateDataset(dataSet.Id, dataSet)
}

func (s *DatasetService) DeleteDataset(dataSetId string) error {
	err := s.validator.ValidateDeletion(dataSetId)
	if err != nil {
		return err
	}
	return s.repo.DeleteDataset(dataSetId)
}

func (s *DatasetService) ListDatasetIds() ([]string, error) {
	return s.repo.ListDatasetIds()
}

func (s *DatasetService) AddBatchToDataset(datasetId string, batch *common.Batch) error {
	err := utils.ValidateId(datasetId)
	if err != nil {
		return fmt.Errorf("error validating dataset id: %v", err)
	}

	exists, err := s.repo.DoesDatasetExist(datasetId)
	if err != nil {
		return fmt.Errorf("error checking if dataset exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("dataset does not exist: %s", datasetId)
	}

	dataset, err := s.GetDataset(datasetId)
	if err != nil {
		return fmt.Errorf("error getting dataset: %v", err)
	}

	if batch.SpecId != dataset.SpecId {
		return fmt.Errorf("batch spec does not match dataset spec")
	}

	err = s.batchService.CreateBatch(datasetId, batch)
	if err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}

	dataset.BatchIds = append(dataset.BatchIds, batch.Id)

	utils.SortStringSlice(dataset.BatchIds)
	err = s.repo.UpdateDataset(datasetId, dataset)
	if err != nil {
		return fmt.Errorf("error updating dataset: %v", err)
	}
	return nil
}

func (s *DatasetService) RemoveBatchFromDataset(dataSetId string, batchId string) error {
	exists, err := s.DoesDatasetExist(dataSetId)
	if err != nil {
		return fmt.Errorf("error checking if dataset exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("dataset does not exist: %s", dataSetId)
	}

	err = s.batchService.DeleteBatch(dataSetId, batchId)
	if err != nil {
		return fmt.Errorf("error deleting batch: %v", err)
	}

	dataset, err := s.GetDataset(dataSetId)
	if err != nil {
		return fmt.Errorf("error getting dataset: %v", err)
	}

	dataset.BatchIds, err = s.batchService.ListBatchIds(dataSetId)
	if err != nil {
		return fmt.Errorf("error listing batch ids: %v", err)
	}

	err = s.repo.UpdateDataset(dataSetId, dataset)
	if err != nil {
		return fmt.Errorf("error updating dataset: %v", err)
	}
	return nil
}

func (s *DatasetService) GetBatchFromDataset(dataSetId, batchId string) (*common.Batch, error) {
	exists, err := s.DoesDatasetExist(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error checking if dataset exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("dataset does not exist: %s", dataSetId)
	}
	batch, err := s.batchService.GetBatch(dataSetId, batchId)
	if err != nil {
		return nil, fmt.Errorf("error getting batch: %v", err)
	}
	return batch, nil
}
