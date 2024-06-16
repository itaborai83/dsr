package datasets

import (
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

type DatasetValidator struct {
	repo     common.DatasetRepo
	specRepo common.SpecRepo
}

func NewDatasetValidator(repo common.DatasetRepo, specRepo common.SpecRepo) *DatasetValidator {
	return &DatasetValidator{repo: repo, specRepo: specRepo}
}

func (d *DatasetValidator) validateSpec(specId string) error {
	err := utils.ValidateId(specId)
	if err != nil {
		return fmt.Errorf("error validating spec id: %v", err)
	}

	exists, err := d.specRepo.DoesSpecExist(specId)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", specId)
	}
	return nil
}

func (d *DatasetValidator) validateDataset(dataset *common.Dataset) error {
	err := utils.ValidateId(dataset.Id)
	if err != nil {
		return fmt.Errorf("error validating dataset id: %v", err)
	}

	if dataset.Name == "" {
		return fmt.Errorf("dataset name is empty")
	}

	err = utils.ValidateId(dataset.SpecId)
	if err != nil {
		return fmt.Errorf("error validating dataset spec id: %v", err)
	}

	if dataset.BatchIds == nil {
		return fmt.Errorf("dataset batch ids is nil")
	}

	seenBatchIds := make(map[string]bool)
	for i, batchId := range dataset.BatchIds {
		err = utils.ValidateId(batchId)
		if err != nil {
			return fmt.Errorf("error validating batch id at index %d: %v", i, err)
		}
		if seenBatchIds[batchId] {
			return fmt.Errorf("duplicate batch id: %s", batchId)
		}
		seenBatchIds[batchId] = true
	}

	err = d.validateSpec(dataset.SpecId)
	if err != nil {
		return fmt.Errorf("error validating dataset spec: %v", err)
	}

	return nil
}

func (d *DatasetValidator) ValidateCreation(dataset *common.Dataset) error {
	err := d.validateDataset(dataset)
	if err != nil {
		return fmt.Errorf("error validating dataset: %v", err)
	}

	exists, err := d.repo.DoesDatasetExist(dataset.Id)
	if err != nil {
		return fmt.Errorf("error checking if dataset exists: %v", err)
	}
	if exists {
		return fmt.Errorf("dataset already exists: %s", dataset.Id)
	}
	return nil
}

func (d *DatasetValidator) ValidateUpdate(datasetId string, dataset *common.Dataset) error {
	if datasetId != dataset.Id {
		return fmt.Errorf("dataset id mismatch")
	}

	err := d.validateDataset(dataset)
	if err != nil {
		return fmt.Errorf("error validating dataset: %v", err)
	}

	exists, err := d.repo.DoesDatasetExist(dataset.Id)
	if err != nil {
		return fmt.Errorf("error checking if dataset exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("dataset does not exist: %s", dataset.Id)
	}

	return nil
}

func (d *DatasetValidator) ValidateDeletion(datasetId string) error {
	err := utils.ValidateId(datasetId)
	if err != nil {
		return fmt.Errorf("error validating dataset id: %v", err)
	}
	exists, err := d.repo.DoesDatasetExist(datasetId)
	if err != nil {
		return fmt.Errorf("error checking if dataset exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("dataset does not exist: %s", datasetId)
	}
	return nil
}
