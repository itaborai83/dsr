package datasets

import (
	"fmt"

	"github.com/itaborai83/dsr/utils"
)

const (
	DATA_SET_INVALID_ID_ERROR       = "invalid id: %v"
	DATA_SET_NAME_EMPTY_ERROR       = "data set name cannot be empty"
	DATA_SET_INVALID_SPEC_ID_ERROR  = "invalid spec id: %v"
	DATA_SET_NIL_BATCH_IDS_ERROR    = "data set batch ids cannot be nil"
	DATA_SET_INVALID_BATCH_ID_ERROR = "invalid batch id at index %d: %v"
)

func ValidateDataSet(dataSet *DataSet) error {
	err := utils.ValidateId(dataSet.Id)
	if err != nil {
		return fmt.Errorf(DATA_SET_INVALID_ID_ERROR, err)
	}

	if dataSet.Name == "" {
		return fmt.Errorf(DATA_SET_NAME_EMPTY_ERROR)
	}

	err = utils.ValidateId(dataSet.SpecId)
	if err != nil {
		return fmt.Errorf(DATA_SET_INVALID_SPEC_ID_ERROR, err)
	}

	if dataSet.BatchIds == nil {
		return fmt.Errorf(DATA_SET_NIL_BATCH_IDS_ERROR)
	}

	for i, batchId := range dataSet.BatchIds {
		err = utils.ValidateId(batchId)
		if err != nil {
			return fmt.Errorf(DATA_SET_INVALID_BATCH_ID_ERROR, i, err)
		}
	}

	return nil
}

/*
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

*/
