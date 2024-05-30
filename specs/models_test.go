package specs

import (
	"testing"
)

func TestColumnSpec_IsValidValue(t *testing.T) {
	tableTests := []struct {
		spec   ColumnSpec
		value  interface{}
		result bool
	}{
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, "value", true},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, 1, true},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, 1.1, true},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, "2020-01-01", true},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, "2020-01-01T00:00:00Z", true},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, true, true},
		{ColumnSpec{Name: "errorType", Type: "INVALID"}, "value", false},
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, 1, false},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, "value", false},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, "value", false},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, 1, false},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, 1, false},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, "value", false},
	}

	for i, tt := range tableTests {
		if tt.spec.IsValidValue(tt.value) != tt.result {
			t.Errorf("Test %d: Expected %v, got %v", i, tt.result, !tt.result)
		}
	}
}

func TestColumnSpec_ConformsTo(t *testing.T) {
	tableTests := []struct {
		spec      ColumnSpec
		tableData interface{}
		result    bool
	}{
		// valid cases
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, map[string]interface{}{"name": []interface{}{"value"}}, true},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, map[string]interface{}{"count": []interface{}{1}}, true},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, map[string]interface{}{"price": []interface{}{1.1}}, true},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, map[string]interface{}{"date": []interface{}{"2020-01-01"}}, true},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, map[string]interface{}{"time": []interface{}{"2020-01-01T00:00:00Z"}}, true},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, map[string]interface{}{"flag": []interface{}{true}}, true},
		// empty data, all valid
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, map[string]interface{}{"name": []interface{}{}}, true},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, map[string]interface{}{"count": []interface{}{}}, true},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, map[string]interface{}{"price": []interface{}{}}, true},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, map[string]interface{}{"date": []interface{}{}}, true},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, map[string]interface{}{"time": []interface{}{}}, true},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, map[string]interface{}{"flag": []interface{}{}}, true},
		// invalid cases - not arrays
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, map[string]interface{}{"name": "value"}, false},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, map[string]interface{}{"count": 1}, false},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, map[string]interface{}{"price": 1.1}, false},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, map[string]interface{}{"date": "2020-01-01"}, false},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, map[string]interface{}{"time": "2020-01-01T00:00:00Z"}, false},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, map[string]interface{}{"flag": true}, false},
		// invalid cases - invalid values
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, map[string]interface{}{"name": []interface{}{1}}, false},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, map[string]interface{}{"count": []interface{}{"value"}}, false},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, map[string]interface{}{"price": []interface{}{"value"}}, false},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, map[string]interface{}{"date": []interface{}{1}}, false},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, map[string]interface{}{"time": []interface{}{1}}, false},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, map[string]interface{}{"flag": []interface{}{"value"}}, false},
		// invalid cases - missing keys
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, map[string]interface{}{}, false},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, map[string]interface{}{}, false},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, map[string]interface{}{}, false},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, map[string]interface{}{}, false},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, map[string]interface{}{}, false},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, map[string]interface{}{}, false},
		// invalid cases - wrong keys
		{ColumnSpec{Name: "name", Type: ColumnTypeString}, map[string]interface{}{"wrong": []interface{}{"value"}}, false},
		{ColumnSpec{Name: "count", Type: ColumnTypeInteger}, map[string]interface{}{"wrong": []interface{}{1}}, false},
		{ColumnSpec{Name: "price", Type: ColumnTypeFloat}, map[string]interface{}{"wrong": []interface{}{1.1}}, false},
		{ColumnSpec{Name: "date", Type: ColumnTypeDate}, map[string]interface{}{"wrong": []interface{}{"2020-01-01"}}, false},
		{ColumnSpec{Name: "time", Type: ColumnTypeDateTime}, map[string]interface{}{"wrong": []interface{}{"2020-01-01T00:00:00Z"}}, false},
		{ColumnSpec{Name: "flag", Type: ColumnTypeBoolean}, map[string]interface{}{"wrong": []interface{}{true}}, false},
	}
	for i, tt := range tableTests {
		if tt.spec.ConformsTo(tt.tableData) != tt.result {
			t.Errorf("Test %d: Expected %v, got %v", i, tt.result, !tt.result)
		}
	}
}

func TestTableSpec_ConformsTo(t *testing.T) {
	spec := TableSpec{
		Name: "name",
		Columns: []ColumnSpec{
			{Name: "name", Type: ColumnTypeString},
			{Name: "count", Type: ColumnTypeInteger},
			{Name: "price", Type: ColumnTypeFloat},
			{Name: "date", Type: ColumnTypeDate},
			{Name: "time", Type: ColumnTypeDateTime},
			{Name: "flag", Type: ColumnTypeBoolean},
		},
		KeyColumns:          []string{"name"},
		ChangeControlColumn: "date",
	}

	tableTests := []struct {
		tableData interface{}
		result    bool
	}{
		// valid case
		{
			map[string]interface{}{
				"name":  []interface{}{"value"},
				"count": []interface{}{1},
				"price": []interface{}{1.1},
				"date":  []interface{}{"2020-01-01"},
				"time":  []interface{}{"2020-01-01T00:00:00Z"},
				"flag":  []interface{}{true},
			},
			true,
		},
		// invalid case - missing key
		{
			map[string]interface{}{
				"count": []interface{}{1},
				"price": []interface{}{1.1},
				"date":  []interface{}{"2020-01-01"},
				"time":  []interface{}{"2020-01-01T00:00:00Z"},
				"flag":  []interface{}{true},
			},
			false,
		},
		// invalid case - wrong key
		{
			map[string]interface{}{
				"wrong": []interface{}{"value"},
				"count": []interface{}{1},
				"price": []interface{}{1.1},
				"date":  []interface{}{"2020-01-01"},
				"time":  []interface{}{"2020-01-01T00:00:00Z"},
				"flag":  []interface{}{true},
			},
			false,
		},
		// invalid case - wrong type
		{
			map[string]interface{}{
				"name":  []interface{}{1},
				"count": []interface{}{"value"},
				"price": []interface{}{"value"},
				"date":  []interface{}{1},
				"time":  []interface{}{1},
				"flag":  []interface{}{"value"},
			},
			false,
		},
	}
	for i, tt := range tableTests {
		if spec.ConformsTo(tt.tableData) != tt.result {
			t.Errorf("Test %d: Expected %v, got %v", i, tt.result, !tt.result)
		}
	}
}
