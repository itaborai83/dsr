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
