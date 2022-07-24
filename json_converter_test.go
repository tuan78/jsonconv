package jsonconv

import (
	"testing"
)

func TestCreateCsvHeader(t *testing.T) {
	jsonArray := []map[string]interface{}{
		{
			"a": true,
			"b": false,
		},
	}
	headers := CreateCsvHeader(jsonArray, nil)
	if len(headers) == 0 {
		t.Fatalf("headers must not be empty")
	}
}
