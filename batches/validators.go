package batches

import (
	"fmt"

	"github.com/itaborai83/dsr/utils"
)

const (
	BATCH_INVALID_ID_ERROR        = "invalid id: %v"
	BATCH_NAME_EMPTY_ERROR        = "batch name cannot be empty"
	BATCH_INVALID_SPEC_ID_ERROR   = "invalid spec id: %v"
	BATCH_INVALID_DATA_SET_ID     = "invalid data set id: %v"
	BATCH_DATA_SET_ID_EMPTY_ERROR = "batch data set id cannot be empty"
	BATCH_DATA_NIL_ERROR          = "batch data cannot be nil"
)

func ValidateBatch(batch *Batch) error {

	err := utils.ValidateId(batch.Id)
	if err != nil {
		return fmt.Errorf(BATCH_INVALID_ID_ERROR, err)
	}
	if batch.Name == "" {
		return fmt.Errorf(BATCH_NAME_EMPTY_ERROR)
	}
	err = utils.ValidateId(batch.SpecId)
	if err != nil {
		return fmt.Errorf(BATCH_INVALID_ID_ERROR, err)
	}

	err = utils.ValidateId(batch.DataSetId)
	if err != nil {
		return fmt.Errorf(BATCH_INVALID_DATA_SET_ID, err)
	}

	if batch.Data == nil {
		return fmt.Errorf("batch data cannot be nil")
	}

	return nil
}
