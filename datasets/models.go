package datasets

const (
	COL_TYPE_STRING   = "STRING"
	COL_TYPE_INTEGER  = "INTEGER"
	COL_TYPE_DATE     = "DATE"
	COL_TYPE_DATETIME = "DATETIME"
)

type ColumnSpec struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
}

type ColumnValues struct {
	Name   string        `json:"column"`
	Values []interface{} `json:"values"`
}

type Batch struct {
	Id        string         `json:"id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	RowCount  int            `json:"row_count"`
	Columns   []ColumnValues `json:"columns"`
}

type BatchSummary struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	RowCount  int    `json:"row_count"`
}

type Dataset struct {
	Id                  string       `json:"id"`
	Name                string       `json:"name"`
	CreatedAt           string       `json:"created_at"`
	UpdatedAt           string       `json:"updated_at"`
	Columns             []ColumnSpec `json:"columns"`
	KeyColumns          []string     `json:"key_columns"`
	ChangeControlColumn string       `json:"change_control_column"`
}
