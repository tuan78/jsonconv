package jsonconv

import (
	"encoding/json"
	"io"
)

// A JsonReader reads and decodes JSON values from an input stream.
type JsonReader struct {
	reader io.Reader
}

// NewJsonReader returns a new JsonReader that reads from r.
func NewJsonReader(r io.Reader) *JsonReader {
	return &JsonReader{
		reader: r,
	}
}

// Read reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (r *JsonReader) Read(v any) error {
	decoder := json.NewDecoder(r.reader)
	for decoder.More() {
		err := decoder.Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}
