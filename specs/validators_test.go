package specs

import (
	"fmt"
	"testing"
)

func TestValidateColumnSpec_ValidColumnSpec(t *testing.T) {
	tableTests := []struct {
		spec   ColumnSpec
		result error
	}{
		// valid cases
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, nil},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, nil},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, nil},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, nil},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, nil},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, nil},
		// invalid case - empty name
		{ColumnSpec{Name: "", Type: ColumnTypeString}, fmt.Errorf(EMPTY_COLUMN_NAME_ERROR)},
		// invalid case - empty type
		{ColumnSpec{Name: "name", Type: ""}, fmt.Errorf(EMPTY_COLUMN_TYPE_ERROR)},
		// invalid case - invalid type
		{ColumnSpec{Name: "name", Type: "INVALID"}, fmt.Errorf(INVALID_COLUMN_TYPE_ERROR, "INVALID")},
	}

	for i, tt := range tableTests {

		err := ValidateColumnSpec(tt.spec)
		if err == nil && tt.result == nil {
			continue
		}
		if err.Error() != tt.result.Error() {
			t.Errorf("Test %d: Expected %v, got %v", i, tt.result, err)
		}
	}
}

func TestValidateTableSpec_ValidTableSpec(t *testing.T) {
	spec := &TableSpec{
		Name: "name",
		Columns: []ColumnSpec{
			{Name: "name", Type: "STRING"},
		},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestValidateTableSpec_EmptyName(t *testing.T) {
	spec := &TableSpec{
		Name: "",
		Columns: []ColumnSpec{
			{Name: "name", Type: "STRING"},
		},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := TABLE_NAME_EMPTY_ERROR
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_ColumnsNil(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             nil,
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := COLUMNS_CANNOT_BE_NIL_ERROR
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_NoColumns(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := NO_COLUMNS_ERROR
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_ColumnValidation(t *testing.T) {
	spec := &TableSpec{
		Name: "name",
		Columns: []ColumnSpec{
			{Name: "", Type: "STRING"},
		},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	columnError := fmt.Errorf(EMPTY_COLUMN_NAME_ERROR)
	expected := fmt.Sprintf(COLUMN_VALIDATION_ERROR, 0, columnError)
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_DuplicateColumn(t *testing.T) {
	spec := &TableSpec{
		Name: "name",
		Columns: []ColumnSpec{
			{Name: "name", Type: "STRING"},
			{Name: "name", Type: "STRING"},
		},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := fmt.Sprintf(DUPLICATE_COLUMN_ERROR, "name")
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_KeyColumnsNil(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          nil,
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := KEY_COLUMNS_NIL_ERROR
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_NoKeyColumns(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          []string{},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := NO_KEY_COLUMNS_ERROR
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_EmptyKeyColumnEntry(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          []string{""},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := fmt.Sprintf(EMPTY_KEY_COLUMN_ENTRY_ERROR, 0)
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_DuplicateKeyColumn(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          []string{"name", "name"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := fmt.Sprintf(DUPLICATE_KEY_COLUMN_ERROR, "name")
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_InvalidKeyColumn(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          []string{"invalid"},
		ChangeControlColumn: "name",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := fmt.Sprintf(INVALID_KEY_COLUMN_ERROR, "invalid")
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_ChangeControlColumnNotInColumns(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "invalid",
	}
	err := ValidateTableSpec(spec)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := fmt.Sprintf(INVALID_CHANGE_CTRL_COLUMN_ERROR, "invalid")
	if err.Error() != expected {
		t.Errorf("Expected %s, got %v", expected, err)
	}
}

func TestValidateTableSpec_ChangeControlColumnEmpty(t *testing.T) {
	spec := &TableSpec{
		Name:                "name",
		Columns:             []ColumnSpec{{Name: "name", Type: "STRING"}},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "",
	}
	err := ValidateTableSpec(spec)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
