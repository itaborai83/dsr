package batches

import (
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

type BatchValidator struct {
	repo     common.BatchRepo
	specRepo common.SpecRepo
}

func NewBatchValidator(repo common.BatchRepo, specRepo common.SpecRepo) *BatchValidator {
	return &BatchValidator{repo: repo, specRepo: specRepo}
}

func (v *BatchValidator) validateBatch(datasetId, batchId string, batch *common.Batch) error {
	if batch == nil {
		return fmt.Errorf("batch is nil")
	}

	if batch.DatasetId != datasetId {
		return fmt.Errorf("batch dataset id mismatch: %s != %s", batch.DatasetId, datasetId)
	}

	if batch.Id != batchId {
		return fmt.Errorf("batch id mismatch: %s != %s", batch.Id, batchId)
	}

	err := utils.ValidateId(batchId)
	if err != nil {
		return fmt.Errorf("error validating batch id: %v", err)
	}

	err = utils.ValidateId(datasetId)
	if err != nil {
		return fmt.Errorf("error validating dataset id: %v", err)
	}

	if batch.Name == "" {
		return fmt.Errorf("batch name is empty")
	}

	err = v.validateSpec(batch.SpecId)
	if err != nil {
		return fmt.Errorf("error validating batch spec id: %v", err)
	}

	return nil
}

func (v *BatchValidator) validateSpec(specId string) error {
	err := utils.ValidateId(specId)
	if err != nil {
		return fmt.Errorf("error validating spec id: %v", err)
	}

	exists, err := v.specRepo.DoesSpecExist(specId)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", specId)
	}
	return nil
}

func (b *BatchValidator) validateAgainstSpec(batch *common.Batch) error {

	exists, err := b.specRepo.DoesSpecExist(batch.SpecId)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", batch.SpecId)
	}

	spec, err := b.specRepo.GetSpec(batch.SpecId)
	if err != nil {
		return fmt.Errorf("error getting spec: %v", err)
	}

	conforms := spec.ConformsTo(batch.Data)
	if !conforms {
		return fmt.Errorf("batch data does not conform to spec: %s", batch.SpecId)
	}
	return nil
}

func (b *BatchValidator) validateColumnLength(batch *common.Batch) error {
	if batch.Data == nil {
		return fmt.Errorf("batch data is nil")
	}

	minColLenght := (1<<31 - 1) // max int size
	maxColLenght := 0

	colLenghts := make(map[string]int)
	// iterate over the keys of the map, i.e., the column names
	for colName := range batch.Data {
		// check if the column is a slice
		colSlice, ok := batch.Data[colName].([]interface{})
		if !ok {
			return fmt.Errorf("column %s is not a slice", colName)
		}
		// store the length of the column
		colLenght := len(colSlice)
		// update the min and max column lengths
		if colLenght < minColLenght {
			minColLenght = colLenght
		}
		if colLenght > maxColLenght {
			maxColLenght = colLenght
		}
		// store the column length
		colLenghts[colName] = len(colSlice)
	}
	if minColLenght != maxColLenght {
		return fmt.Errorf("column lengths are not equal: %v", colLenghts)
	}
	batch.RecordCount = minColLenght
	return nil
}

func (v *BatchValidator) ValidateCreation(datasetId string, batch *common.Batch) error {
	err := v.validateBatch(datasetId, batch.Id, batch)
	if err != nil {
		return fmt.Errorf("error validating batch: %v", err)
	}

	err = v.validateAgainstSpec(batch)
	if err != nil {
		return fmt.Errorf("error validating batch against spec: %v", err)
	}

	err = v.validateColumnLength(batch)
	if err != nil {
		return fmt.Errorf("error validating column length: %v", err)
	}
	return nil
}

func (v *BatchValidator) ValidateDeletion(datasetId, batchId string) error {
	err := utils.ValidateId(batchId)
	if err != nil {
		return fmt.Errorf("error validating batch id: %v", err)
	}

	err = utils.ValidateId(datasetId)
	if err != nil {
		return fmt.Errorf("error validating dataset id: %v", err)
	}

	exists, err := v.repo.DoesBatchExist(datasetId, batchId)
	if err != nil {
		return fmt.Errorf("error checking if batch exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("batch does not exist: %s", batchId)
	}
	return nil
}
