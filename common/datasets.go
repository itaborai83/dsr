package common

type Dataset struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	SpecId    string   `json:"spec_id"`
	BatchIds  []string `json:"batch_ids"`
}
