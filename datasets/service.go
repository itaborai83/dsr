package datasets

import (
	"encoding/json"
	"fmt"

	"github.com/itaborai83/dsr/batches"
	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/specs"
	"github.com/itaborai83/dsr/utils"
)

type DataSetService struct {
	repo         common.Repo
	SpecService  *specs.SpecService
	BatchService *batches.BatchService
}

func RegisterServices() error {
	conf := config.GetConfig()
	err := utils.EnsureFolderExists(conf.DataDir, conf.DataSetsDir)
	if err != nil {
		return fmt.Errorf("error ensuring datasets directory exists: %v", err)
	}

	datasetsDir := conf.DataDir + "/" + conf.DataSetsDir
	repo, err := common.NewFileRepo(datasetsDir, conf.DataSetEntry)
	if err != nil {
		return fmt.Errorf("error creating datasets repo: %v", err)
	}

	specService := conf.SpecService.(*specs.SpecService)
	if specService == nil {
		return fmt.Errorf("error getting previously registered spec service")
	}

	batchService := conf.BatchService.(*batches.BatchService)
	if batchService == nil {
		return fmt.Errorf("error getting previously registered batch service")
	}

	service := NewDataSetService(repo, specService, batchService)
	conf.DataSetsRepo = repo
	conf.DataSetService = service
	return nil
}

func NewDataSetService(repo common.Repo, specService *specs.SpecService, batchService *batches.BatchService) *DataSetService {
	return &DataSetService{
		repo:         repo,
		SpecService:  specService,
		BatchService: batchService,
	}
}

func (d *DataSetService) saveDataSet(dataSet *DataSet) error {
	now := utils.GetNow()
	create := dataSet.CreatedAt == ""

	if create {
		dataSet.CreatedAt = now
	}
	dataSet.UpdatedAt = now

	data, err := json.Marshal(dataSet)
	if err != nil {
		return fmt.Errorf("error marshalling data set: %v", err)
	}
	if create {
		err = d.repo.CreateEntry(dataSet.Id, data)
	} else {
		err = d.repo.UpdateEntry(dataSet.Id, data)
	}

	if err != nil {
		if create {
			return fmt.Errorf("error creating data set: %v", err)
		} else {
			return fmt.Errorf("error updating data set: %v", err)
		}
	}
	return nil
}

func (d *DataSetService) loadDataSet(dataSetId string) (*DataSet, error) {
	data, err := d.repo.GetEntry(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error getting data set: %v", err)
	}
	var dataSet DataSet
	err = json.Unmarshal(data, &dataSet)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data set: %v", err)
	}
	return &dataSet, nil
}

func (d *DataSetService) DoesDataSetExist(dataSetId string) (bool, error) {
	return d.repo.DoesEntryExist(dataSetId)
}

func (d *DataSetService) validateSpecExists(dataSet DataSet) error {
	exists, err := d.SpecService.DoesSpecExist(dataSet.SpecId)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", dataSet.SpecId)
	}
	return nil
}

func (d *DataSetService) validateGetDataSet(dataSetId string) error {
	err := utils.ValidateId(dataSetId)
	if err != nil {
		return fmt.Errorf("error validating data set id: %v", err)
	}

	exists, err := d.DoesDataSetExist(dataSetId)
	if err != nil {
		return fmt.Errorf("error checking if data set exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("data set does not exist: %s", dataSetId)
	}

	return nil

}

func (d *DataSetService) GetDataSet(dataSetId string) (*DataSet, error) {

	err := d.validateGetDataSet(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error validating data set: %v", err)
	}

	dataSet, err := d.loadDataSet(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error loading data set: %v", err)
	}
	return dataSet, nil
}

func (d *DataSetService) validateCreation(dataSet *DataSet) error {
	err := ValidateDataSet(dataSet)
	if err != nil {
		return fmt.Errorf("error validating data set: %v", err)
	}

	exists, err := d.DoesDataSetExist(dataSet.Id)
	if err != nil {
		return fmt.Errorf("error checking if data set exists: %v", err)
	}
	if exists {
		return fmt.Errorf("data set already exists: %s", dataSet.Id)
	}

	err = d.validateSpecExists(*dataSet)
	if err != nil {
		return fmt.Errorf("error validating spec: %v", err)
	}

	if len(dataSet.BatchIds) > 0 {
		return fmt.Errorf("data set cannot have batch ids on creation")
	}

	return nil
}

func (d *DataSetService) CreateDataSet(dataSet *DataSet) error {

	err := d.validateCreation(dataSet)
	if err != nil {
		return fmt.Errorf("error validating data set: %v", err)
	}

	err = d.saveDataSet(dataSet)
	if err != nil {
		return fmt.Errorf("error saving data set: %v", err)
	}
	return nil
}

func (d *DataSetService) validateDeletion(dataSetId string) error {
	err := utils.ValidateId(dataSetId)
	if err != nil {
		return fmt.Errorf("error validating data set id: %v", err)
	}

	exists, err := d.DoesDataSetExist(dataSetId)
	if err != nil {
		return fmt.Errorf("error checking if data set exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("data set does not exist: %s", dataSetId)
	}

	return nil
}

func (d *DataSetService) DeleteDataSet(dataSetId string) error {
	err := d.validateDeletion(dataSetId)
	if err != nil {
		return fmt.Errorf("error validating data set: %v", err)
	}

	dataSet, err := d.GetDataSet(dataSetId)
	if err != nil {
		return fmt.Errorf("error getting data set: %v", err)
	}

	// remove each batch
	copyBatchIds := make([]string, len(dataSet.BatchIds))
	copy(copyBatchIds, dataSet.BatchIds)
	for _, batchId := range copyBatchIds {
		err = d.BatchService.DeleteBatch(batchId)
		if err != nil {
			return fmt.Errorf("error deleting batch: %v", err)
		}

		// remove batch id from data set and save it
		dataSet.BatchIds = utils.RemoveStringFromList(dataSet.BatchIds, batchId)
		err := d.saveDataSet(dataSet)
		if err != nil {
			return fmt.Errorf("error saving data set: %v", err)
		}
	}

	err = d.repo.DeleteEntry(dataSetId)
	if err != nil {
		return fmt.Errorf("error deleting data set: %v", err)
	}
	return nil
}

func (d *DataSetService) ListDataSetIds() ([]string, error) {
	ids, err := d.repo.ListEntryIds()
	if err != nil {
		return nil, fmt.Errorf("error listing data set ids: %v", err)
	}
	return ids, nil
}

func (d *DataSetService) ListDataSets() ([]DataSet, error) {
	ids, err := d.ListDataSetIds()
	if err != nil {
		return nil, err
	}

	dataSets := make([]DataSet, len(ids))
	for i, id := range ids {
		dataSet, err := d.GetDataSet(id)
		if err != nil {
			return nil, err
		}
		dataSets[i] = *dataSet
	}
	return dataSets, nil
}

func (d *DataSetService) validateAddBatchToDataSet(dataSetId string, batch *batches.Batch) (*DataSet, error) {
	err := utils.ValidateId(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error validating data set id: %v", err)
	}

	err = utils.ValidateId(batch.Id)
	if err != nil {
		return nil, fmt.Errorf("error validating batch id: %v", err)
	}

	exists, err := d.DoesDataSetExist(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error checking if data set exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("data set does not exist: %s", dataSetId)
	}

	exists, err = d.BatchService.DoesBatchExist(batch.Id)
	if err != nil {
		return nil, fmt.Errorf("error checking if batch exists: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("batch already exists: %s", batch.Id)
	}

	// get data set
	dataSet, err := d.GetDataSet(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error getting data set: %v", err)
	}

	// check spec id
	if batch.SpecId != dataSet.SpecId {
		return nil, fmt.Errorf("batch spec id mismatch")
	}

	// check if spec exists
	err = d.validateSpecExists(*dataSet)
	if err != nil {
		return nil, fmt.Errorf("error validating spec: %v", err)
	}

	// check if batch is already in data set
	if utils.StringInList(dataSet.BatchIds, batch.Id) {
		return nil, fmt.Errorf("batch already in data set: %s", batch.Id)
	}
	return dataSet, nil
}

func (d *DataSetService) AddBatchToDataSet(dataSetId string, batch *batches.Batch) error {
	dataSet, err := d.validateAddBatchToDataSet(dataSetId, batch)
	if err != nil {
		return fmt.Errorf("error validating data set: %v", err)
	}

	err = d.BatchService.CreateBatch(batch)
	if err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}

	dataSet.BatchIds = append(dataSet.BatchIds, batch.Id)

	err = d.saveDataSet(dataSet)
	if err != nil {
		return fmt.Errorf("error saving data set: %v", err)
	}
	return nil
}

func (d *DataSetService) validateRemoveBatchFromDataSet(dataSetId, batchId string) (*DataSet, error) {
	err := utils.ValidateId(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error validating data set id: %v", err)
	}

	err = utils.ValidateId(batchId)
	if err != nil {
		return nil, fmt.Errorf("error validating batch id: %v", err)
	}

	exists, err := d.DoesDataSetExist(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error checking if data set exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("data set does not exist: %s", dataSetId)
	}

	exists, err = d.BatchService.DoesBatchExist(batchId)
	if err != nil {
		return nil, fmt.Errorf("error checking if batch exists: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("batch already exists: %s", batchId)
	}

	// get data set
	dataSet, err := d.GetDataSet(dataSetId)
	if err != nil {
		return nil, fmt.Errorf("error getting data set: %v", err)
	}
	// get batch
	batch, err := d.BatchService.GetBatch(batchId)
	if err != nil {
		return nil, fmt.Errorf("error getting batch: %v", err)
	}
	// check spec id
	if batch.SpecId != dataSet.SpecId {
		return nil, fmt.Errorf("batch spec id mismatch")
	}

	// check if spec exists
	err = d.validateSpecExists(*dataSet)
	if err != nil {
		return nil, fmt.Errorf("error validating spec: %v", err)
	}

	// check if batch is already in data set
	if utils.StringInList(dataSet.BatchIds, batchId) {
		return nil, fmt.Errorf("batch already in data set: %s", batchId)
	}
	return dataSet, nil
}

func (d *DataSetService) RemoveBatchFromDataSet(dataSetId string, batchId string) error {
	dataSet, err := d.validateRemoveBatchFromDataSet(dataSetId, batchId)
	if err != nil {
		return fmt.Errorf("error validating data set: %v", err)
	}

	err = d.BatchService.DeleteBatch(batchId)
	if err != nil {
		return fmt.Errorf("error deleting batch: %v", err)
	}

	dataSet.BatchIds = utils.RemoveStringFromList(dataSet.BatchIds, batchId)

	err = d.saveDataSet(dataSet)
	if err != nil {
		return fmt.Errorf("error saving data set: %v", err)
	}
	return nil
}

func (d *DataSetService) validateGetBatchFromDataSet(dataSetId, batchId string) error {
	err := utils.ValidateId(dataSetId)
	if err != nil {
		return fmt.Errorf("error validating data set id: %v", err)
	}

	err = utils.ValidateId(batchId)
	if err != nil {
		return fmt.Errorf("error validating batch id: %v", err)
	}

	dataSet, err := d.GetDataSet(dataSetId)
	if err != nil {
		return fmt.Errorf("error getting data set: %v", err)
	}

	if !utils.StringInList(dataSet.BatchIds, batchId) {
		return fmt.Errorf("batch not in data set: %s", batchId)
	}

	return nil
}

func (d *DataSetService) GetBatchFromDataSet(dataSetId, batchId string) (*batches.Batch, error) {
	err := d.validateGetBatchFromDataSet(dataSetId, batchId)
	if err != nil {
		return nil, fmt.Errorf("error validating data set: %v", err)
	}
	batch, err := d.BatchService.GetBatch(batchId)
	if err != nil {
		return nil, fmt.Errorf("error getting batch: %v", err)
	}
	return batch, nil
}
