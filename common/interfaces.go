package common

type Repo interface {
	HasParentRepo() bool
	GetParentRepo() Repo
	DoesEntryExist(id []string) (bool, error)
	GetEntry(id []string) ([]byte, error)
	CreateEntry(id []string, data []byte) error
	UpdateEntry(id []string, data []byte) error
	DeleteEntry(id []string) error
	ListEntryIds(parentId []string) ([][]string, error)
}

type SpecRepo interface {
	DoesSpecExist(id string) (bool, error)
	GetSpec(id string) (*TableSpec, error)
	CreateSpec(id string, spec *TableSpec) error
	UpdateSpec(id string, spec *TableSpec) error
	DeleteSpec(id string) error
	ListSpecIds() ([]string, error)
}

type SpecValidator interface {
	ValidateCreation(spec *TableSpec) error
	ValidateUpdate(specId string, spec *TableSpec) error
	ValidateDeletion(specId string) error
}

type SpecService interface {
	DoesSpecExist(specId string) (bool, error)
	GetSpec(specId string) (*TableSpec, error)
	CreateSpec(spec *TableSpec) error
	UpdateSpec(specId string, spec *TableSpec) error
	DeleteSpec(specId string) error
	ListSpecIds() ([]string, error)
	GetAllSpecs() ([]TableSpec, error)
}

type DatasetRepo interface {
	DoesDatasetExist(id string) (bool, error)
	GetDataset(id string) (*Dataset, error)
	CreateDataset(id string, dataset *Dataset) error
	UpdateDataset(id string, dataset *Dataset) error
	DeleteDataset(id string) error
	ListDatasetIds() ([]string, error)
}

type DatasetValidator interface {
	ValidateCreation(dataset *Dataset) error
	ValidateUpdate(datasetId string, dataset *Dataset) error
	ValidateDeletion(datasetId string) error
}

type DatasetService interface {
	DoesDatasetExist(dataSetId string) (bool, error)
	GetDataset(dataSetId string) (*Dataset, error)
	CreateDataset(dataSet *Dataset) error
	DeleteDataset(dataSetId string) error
	ListDatasetIds() ([]string, error)
	AddBatchToDataset(dataSetId string, batch *Batch) error
	RemoveBatchFromDataset(dataSetId string, batchId string) error
	GetBatchFromDataset(dataSetId, batchId string) (*Batch, error)
}

type BatchRepo interface {
	DoesBatchExist(datasetId, id string) (bool, error)
	GetBatch(datasetId, id string) (*Batch, error)
	CreateBatch(datasetId, id string, batch *Batch) error
	UpdateBatch(datasetId, id string, batch *Batch) error
	DeleteBatch(datasetId, id string) error
	ListBatchIds(datasetId string) ([]string, error)
}

type BatchValidator interface {
	ValidateCreation(datasetId string, batch *Batch) error
	ValidateDeletion(datasetId, batchId string) error
}

type BatchService interface {
	DoesBatchExist(datasetId, batchId string) (bool, error)
	GetBatch(datasetId, batchId string) (*Batch, error)
	CreateBatch(datasetId string, batch *Batch) error
	DeleteBatch(datasetId, batchId string) error
	ListBatchIds(datasetId string) ([]string, error)
}
