package jsonconv

import (
	"strings"
	"testing"
)

func TestJsonReader_InvalidJson(t *testing.T) {
	// Prepare
	raw := `"id": "b042ab5c-ca73-4460-b739-96410ea9d3a6" }`
	obj := make(JsonObject)
	re := NewJsonReader(strings.NewReader(raw))

	// Process
	err := re.Read(&obj)

	// Check
	if err == nil {
		t.Fatalf("Should throw an error for invalid JSON")
	}
}

func TestJsonReader_JsonObject(t *testing.T) {
	// Prepare
	raw := `
		{
			"id": "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user": "Jon Doe",
			"score": "-100",
			"is active": "false"
		}`
	obj := make(JsonObject)
	re := NewJsonReader(strings.NewReader(raw))

	// Process
	err := re.Read(&obj)
	if err != nil {
		t.Fatalf("failed to read json, err: %v", err)
	}

	// Check
	if obj["id"] != "b042ab5c-ca73-4460-b739-96410ea9d3a6" ||
		obj["user"] != "Jon Doe" ||
		obj["score"] != "-100" ||
		obj["is active"] != "false" {
		t.Fatalf("failed to read json object")
	}
}

func TestJsonReader_JsonArray(t *testing.T) {
	// Prepare
	raw := `
	[
		{
			"id": "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
			"user": "Tuấn",
			"score": "1.5",
			"is active": "true"
		},
		{
			"id": "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user": "Jon Doe",
			"score": "-100",
			"is active": "false"
		},
		{
			"id": "4e01b638-44e5-4079-8043-baabbff21cc8",
			"user": "高橋",
			"score": "100000000000000000000000",
			"is active": "true"
		}
	]`
	arr := make(JsonArray, 0)
	re := NewJsonReader(strings.NewReader(raw))

	// Process
	err := re.Read(&arr)
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}

	// Check
	if len(arr) != 3 {
		t.Fatalf("failed to read json array")
	}
	if arr[0]["id"] != "ce06f5b1-5721-42c0-91e1-9f72a09c250a" ||
		arr[0]["user"] != "Tuấn" ||
		arr[0]["score"] != "1.5" ||
		arr[0]["is active"] != "true" {
		t.Fatalf("failed to read json array")
	}
	if arr[1]["id"] != "b042ab5c-ca73-4460-b739-96410ea9d3a6" ||
		arr[1]["user"] != "Jon Doe" ||
		arr[1]["score"] != "-100" ||
		arr[1]["is active"] != "false" {
		t.Fatalf("failed to read json array")
	}
	if arr[2]["id"] != "4e01b638-44e5-4079-8043-baabbff21cc8" ||
		arr[2]["user"] != "高橋" ||
		arr[2]["score"] != "100000000000000000000000" ||
		arr[2]["is active"] != "true" {
		t.Fatalf("failed to read json array")
	}
}

func TestJsonReader_JsonArray_NewlineDelimited(t *testing.T) {
	// Prepare
	raw := `
		{"id": "ce06f5b1-5721-42c0-91e1-9f72a09c250a","user": "Tuấn","score": "1.5","is active": "true"}
		{"id": "b042ab5c-ca73-4460-b739-96410ea9d3a6","user": "Jon Doe","score": "-100","is active": "false"}
		{"id": "4e01b638-44e5-4079-8043-baabbff21cc8","user": "高橋","score": "100000000000000000000000","is active": "true"}`
	arr := make(JsonArray, 0)
	re := NewJsonReader(strings.NewReader(raw))

	// Process
	err := re.Read(&arr)
	if err != nil {
		t.Fatalf("failed to read file %v", err)
	}

	// Check
	if len(arr) != 3 {
		t.Fatalf("failed to read json array")
	}
	if arr[0]["id"] != "ce06f5b1-5721-42c0-91e1-9f72a09c250a" ||
		arr[0]["user"] != "Tuấn" ||
		arr[0]["score"] != "1.5" ||
		arr[0]["is active"] != "true" {
		t.Fatalf("failed to read json array")
	}
	if arr[1]["id"] != "b042ab5c-ca73-4460-b739-96410ea9d3a6" ||
		arr[1]["user"] != "Jon Doe" ||
		arr[1]["score"] != "-100" ||
		arr[1]["is active"] != "false" {
		t.Fatalf("failed to read json array")
	}
	if arr[2]["id"] != "4e01b638-44e5-4079-8043-baabbff21cc8" ||
		arr[2]["user"] != "高橋" ||
		arr[2]["score"] != "100000000000000000000000" ||
		arr[2]["is active"] != "true" {
		t.Fatalf("failed to read json array")
	}
}
