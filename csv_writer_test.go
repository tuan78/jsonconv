package jsonconv

import (
	"bytes"
	"testing"
)

func TestCsvWriter_InvalidDelimiter(t *testing.T) {
	// Prepare
	data := CsvData{
		{
			"id", "user", "score", "is active",
		},
	}
	buf := &bytes.Buffer{}
	wr := NewCsvWriter(buf)
	wr.Delimiter = NewDelimiter('\n')

	// Process
	err := wr.Write(data)

	// Check
	if err == nil {
		t.Fatalf("Should throw an error for invalid delimiter")
	}
}

func TestCsvWriter(t *testing.T) {
	// Prepare
	data := CsvData{
		{
			"id", "user", "score", "is active",
		},
		{
			"ce06f5b1-5721-42c0-91e1-9f72a09c250a", "Tuấn", "1.5", "true",
		},
		{
			"b042ab5c-ca73-4460-b739-96410ea9d3a6", "Jon Doe", "-100", "false",
		},
		{
			"4e01b638-44e5-4079-8043-baabbff21cc8", "高橋", "100000000000000000000000", "true",
		},
		{
			"6f0d6265-545c-4366-a78b-4f80c337aa69", "김슬기", "1234567890", "true",
		},
		{
			"3fbae214-006d-4ac5-9eea-76c5d611f54a", "Comma,", "0", "false",
		},
	}
	buf := &bytes.Buffer{}
	wr := NewCsvWriter(buf)
	wr.Delimiter = NewDelimiter('|')

	// Process
	err := wr.Write(data)
	if err != nil {
		t.Fatalf("failed to write csv, err: %v", err)
	}

	// Check
	s := buf.String()
	expect := `id|user|score|is active
ce06f5b1-5721-42c0-91e1-9f72a09c250a|Tuấn|1.5|true
b042ab5c-ca73-4460-b739-96410ea9d3a6|Jon Doe|-100|false
4e01b638-44e5-4079-8043-baabbff21cc8|高橋|100000000000000000000000|true
6f0d6265-545c-4366-a78b-4f80c337aa69|김슬기|1234567890|true
3fbae214-006d-4ac5-9eea-76c5d611f54a|Comma,|0|false
`
	if s == "" {
		t.Fatalf("failed to write csv to byte buffer")
	}
	if s != expect {
		t.Fatalf("csv output is not correct")
	}
}
