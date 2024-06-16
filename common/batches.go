package common

type Batch struct {
	Id           string                 `json:"id"`
	Name         string                 `json:"name"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
	SpecId       string                 `json:"spec_id"`
	DatasetId    string                 `json:"data_set_id"`
	RecordCount  int                    `json:"record_count"`
	RecordHashes []uint64               `json:"record_hashes"`
	Data         map[string]interface{} `json:"data"`
}
