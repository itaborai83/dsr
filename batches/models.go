package batches

type Batch struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	SpecId    string                 `json:"spec_id"`
	DataSetId string                 `json:"data_set_id"`
	Data      map[string]interface{} `json:"data"`
}
