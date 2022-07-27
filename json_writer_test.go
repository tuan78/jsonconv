package jsonconv

import (
	"bytes"
	"testing"
)

func TestJsonWriter_JsonObject(t *testing.T) {
	// Prepare
	data := JsonObject{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"special1":  "&",
		"special2":  "<",
		"special3":  ">",
		"special4":  "\u0026",
		"special5":  "\u003c",
		"special6":  "\u003e",
	}
	buf := &bytes.Buffer{}
	wr := NewJsonWriter(buf)

	// Process
	err := wr.Write(data)
	if err != nil {
		t.Fatalf("failed to write json, err: %v", err)
	}

	// Check
	s := buf.String()
	expect := `{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"score":-100,"special1":"\u0026","special2":"\u003c","special3":"\u003e","special4":"\u0026","special5":"\u003c","special6":"\u003e","user":"Jon Doe"}
`
	if s == "" {
		t.Fatalf("failed to write json to byte buffer")
	}
	if s != expect {
		t.Fatalf("json output is not correct")
	}
}

func TestJsonWriter_JsonArray(t *testing.T) {
	// Prepare
	data := JsonArray{
		{
			"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user":      "Jon Doe",
			"score":     -100,
			"is active": false,
			"special1":  "&",
			"special2":  "<",
			"special3":  ">",
			"special4":  "\u0026",
			"special5":  "\u003c",
			"special6":  "\u003e",
		},
		{
			"id":        "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
			"user":      "Tuấn",
			"score":     1.5,
			"is active": true,
		},
		{
			"id":        "4e01b638-44e5-4079-8043-baabbff21cc8",
			"user":      "高橋",
			"score":     100000000000000000,
			"is active": true,
		},
	}
	buf := &bytes.Buffer{}
	wr := NewJsonWriter(buf)

	// Process
	err := wr.Write(data)
	if err != nil {
		t.Fatalf("failed to write json, err: %v", err)
	}

	// Check
	s := buf.String()
	expected := `[{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"score":-100,"special1":"\u0026","special2":"\u003c","special3":"\u003e","special4":"\u0026","special5":"\u003c","special6":"\u003e","user":"Jon Doe"},{"id":"ce06f5b1-5721-42c0-91e1-9f72a09c250a","is active":true,"score":1.5,"user":"Tuấn"},{"id":"4e01b638-44e5-4079-8043-baabbff21cc8","is active":true,"score":100000000000000000,"user":"高橋"}]
`
	if s == "" {
		t.Fatalf("failed to write json to byte buffer")
	}
	if s != expected {
		t.Fatalf("json output is not correct")
	}
}
