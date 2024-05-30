package common

type Repo interface {
	DoesEntryExist(entryId string) (bool, error)
	GetEntry(entryId string) ([]byte, error)
	CreateEntry(entryId string, data []byte) error
	UpdateEntry(entryId string, data []byte) error
	DeleteEntry(entryId string) error
	ListEntryIds() ([]string, error)
}
