package jsonconv

import (
	"path/filepath"
	"testing"
)

func TestJsonReaderFromString(t *testing.T) {
	rawData := `
		{
			"a": "test",
			"b": "test"
		}`
	jsonObject := make(JsonObject)
	reader := NewJsonReaderFromString(rawData)
	err := reader.Read(&jsonObject)
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}
	if _, exist := jsonObject["a"]; !exist {
		t.Fatalf("failed to read json object")
	}
}

func TestJsonReaderFromFile_JsonObject(t *testing.T) {
	jsonObject := make(JsonObject)
	reader, err := NewJsonReaderFromFile(filepath.Join("tests", "read-json-object.json"))
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}
	err = reader.Read(&jsonObject)
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}
	if _, exist := jsonObject["a"]; !exist {
		t.Fatalf("failed to read json object")
	}
}

func TestJsonReaderFromFile_JsonArray(t *testing.T) {
	jsonArray := make(JsonArray, 0)
	reader, err := NewJsonReaderFromFile(filepath.Join("tests", "read-json-array.json"))
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}
	err = reader.Read(&jsonArray)
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}
	if len(jsonArray) == 0 {
		t.Fatalf("failed to read json array")
	}
}
