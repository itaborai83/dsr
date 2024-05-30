package specs

const (
	ColumnTypeString   = "STRING"
	ColumnTypeInteger  = "INTEGER"
	ColumnTypeFloat    = "FLOAT"
	ColumnTypeDate     = "DATE"
	ColumnTypeDateTime = "DATETIME"
	ColumnTypeBoolean  = "BOOLEAN"
)

type ColumnSpec struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TableSpec struct {
	Id                  string       `json:"id"`
	Name                string       `json:"name"`
	CreatedAt           string       `json:"created_at"`
	UpdatedAt           string       `json:"updated_at"`
	Columns             []ColumnSpec `json:"columns"`
	KeyColumns          []string     `json:"key_columns"`
	ChangeControlColumn string       `json:"change_control_column"`
}

func (c ColumnSpec) IsValidValue(value interface{}) bool {
	switch c.Type {
	case ColumnTypeString:
		_, ok := value.(string)
		return ok
	case ColumnTypeInteger:
		_, ok := value.(int)
		return ok
	case ColumnTypeFloat:
		_, ok := value.(float64)
		return ok
	case ColumnTypeDate:
		_, ok := value.(string)
		return ok
	case ColumnTypeDateTime:
		_, ok := value.(string)
		return ok
	case ColumnTypeBoolean:
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

func (c ColumnSpec) ConformsTo(tableData interface{}) bool {
	// see if table data is a map
	data, ok := tableData.(map[string]interface{})
	if !ok {
		return false
	}
	// see if column name is in the map
	value, ok := data[c.Name]
	if !ok {
		return false
	}
	// see if value is an array of interface{}
	_, ok = value.([]interface{})
	if !ok {
		return false
	}
	arrayOfValues := value.([]interface{})
	// see if value conforms to the column type
	for _, v := range arrayOfValues {
		if !c.IsValidValue(v) {
			return false
		}
	}
	return ok
}

func (t TableSpec) GetColumn(name string) *ColumnSpec {
	for _, col := range t.Columns {
		if col.Name == name {
			return &col
		}
	}
	return nil
}

func (t TableSpec) ConformsTo(data interface{}) bool {
	// see if table data is a map
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	// see if all the columns are present
	for _, col := range t.Columns {
		if _, ok := dataMap[col.Name]; !ok {
			return false
		}
	}
	// see if all the columns conform to the column types
	for _, col := range t.Columns {
		if !col.ConformsTo(data) {
			return false
		}
	}
	return true
}
