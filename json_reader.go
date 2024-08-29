package jsonconv

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

// A JsonReader reads and decodes JSON values from an input stream.
type JsonReader struct {
	reader io.ReadSeeker
}

// NewJsonReader returns a new JsonReader that reads from r.
func NewJsonReader(r io.ReadSeeker) *JsonReader {
	return &JsonReader{
		reader: r,
	}
}

// Read reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (r *JsonReader) Read(v interface{}) error {
	decoder := json.NewDecoder(r.reader)
	err := decoder.Decode(v)
	if err != nil {
		if !strings.Contains(err.Error(), "cannot unmarshal object into Go value of type") {
			return err
		}

		// The JSON data is a valid JSON array. However, we will try to decode it line by line.
		// Reset reader and decoder.
		r.reader.Seek(0, io.SeekStart)
		decoder = json.NewDecoder(r.reader)

		// Check if v is a pointer to a slice or an array. If not, return an error.
		refval := reflect.ValueOf(v)
		if refval.Kind() == reflect.Pointer {
			refval = refval.Elem()
		}
		if refval.Kind() != reflect.Slice && refval.Kind() != reflect.Array {
			return err
		}

		// Decode JSON array into v.
		for {
			obj := reflect.New(refval.Type().Elem())
			objInterface := obj.Interface()
			err := decoder.Decode(objInterface)
			if err != nil {
				if err == io.EOF {
					reflect.ValueOf(v).Elem().Set(refval)
					return nil
				}
				return err
			}
			refval.Set(reflect.Append(refval, reflect.ValueOf(objInterface).Elem()))
		}
	}
	return nil
}
