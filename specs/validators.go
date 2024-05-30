package specs

import (
	"fmt"

	"github.com/itaborai83/dsr/utils"
)

const (
	TABLE_INVALID_ID                 = "invalid id: %v"
	EMPTY_COLUMN_NAME_ERROR          = "column name cannot be empty"
	EMPTY_COLUMN_TYPE_ERROR          = "column type cannot be empty"
	INVALID_COLUMN_TYPE_ERROR        = "invalid column type: %s"
	TABLE_NAME_EMPTY_ERROR           = "table name cannot be empty"
	COLUMNS_CANNOT_BE_NIL_ERROR      = "columns cannot be nil"
	NO_COLUMNS_ERROR                 = "table must have at least one column"
	COLUMN_VALIDATION_ERROR          = "invalid column at index %d: %s"
	DUPLICATE_COLUMN_ERROR           = "duplicate column name: %s"
	KEY_COLUMNS_NIL_ERROR            = "key columns cannot be nil"
	NO_KEY_COLUMNS_ERROR             = "table must have at least one key column"
	EMPTY_KEY_COLUMN_ENTRY_ERROR     = "key column entry at index %d cannot be empty"
	DUPLICATE_KEY_COLUMN_ERROR       = "duplicate key column: %s"
	INVALID_KEY_COLUMN_ERROR         = "key column %s not found in columns"
	INVALID_CHANGE_CTRL_COLUMN_ERROR = "change control column %s not found in columns"
)

func isValidColumnType(t string) bool {
	return t == ColumnTypeString ||
		t == ColumnTypeInteger ||
		t == ColumnTypeFloat ||
		t == ColumnTypeDate ||
		t == ColumnTypeDateTime ||
		t == ColumnTypeBoolean
}

func ValidateColumnSpec(spec ColumnSpec) error {
	if spec.Name == "" {
		return fmt.Errorf(EMPTY_COLUMN_NAME_ERROR)
	}
	if spec.Type == "" {
		return fmt.Errorf(EMPTY_COLUMN_TYPE_ERROR)
	}
	if !isValidColumnType(spec.Type) {
		return fmt.Errorf(INVALID_COLUMN_TYPE_ERROR, spec.Type)
	}

	return nil
}

func ValidateTableSpec(spec *TableSpec) error {
	err := utils.ValidateId(spec.Id)
	if err != nil {
		return fmt.Errorf(TABLE_INVALID_ID, err)
	}
	if spec.Name == "" {
		return fmt.Errorf(TABLE_NAME_EMPTY_ERROR)
	}
	if spec.Columns == nil {
		return fmt.Errorf(COLUMNS_CANNOT_BE_NIL_ERROR)
	}
	if len(spec.Columns) == 0 {
		return fmt.Errorf(NO_COLUMNS_ERROR)
	}
	seenColumns := make(map[string]bool)
	for i, col := range spec.Columns {
		if err := ValidateColumnSpec(col); err != nil {
			return fmt.Errorf(COLUMN_VALIDATION_ERROR, i, err)
		}
		if seenColumns[col.Name] {
			return fmt.Errorf(DUPLICATE_COLUMN_ERROR, col.Name)
		}
		seenColumns[col.Name] = true
	}

	if spec.KeyColumns == nil {
		return fmt.Errorf(KEY_COLUMNS_NIL_ERROR)
	}
	if len(spec.KeyColumns) == 0 {
		return fmt.Errorf(NO_KEY_COLUMNS_ERROR)
	}
	seenKeyColumns := make(map[string]bool)
	for i, key := range spec.KeyColumns {
		if key == "" {
			return fmt.Errorf(EMPTY_KEY_COLUMN_ENTRY_ERROR, i)
		}
		if seenKeyColumns[key] {
			return fmt.Errorf(DUPLICATE_KEY_COLUMN_ERROR, key)
		}
		seenKeyColumns[key] = true
		if !seenColumns[key] {
			return fmt.Errorf(INVALID_KEY_COLUMN_ERROR, key)
		}
	}

	if spec.ChangeControlColumn != "" {
		if !seenColumns[spec.ChangeControlColumn] {
			return fmt.Errorf(INVALID_CHANGE_CTRL_COLUMN_ERROR, spec.ChangeControlColumn)
		}
	}
	return nil
}
