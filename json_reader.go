package jsonconv

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type JsonReader struct {
	reader io.Reader
	closer io.Closer
}

func NewJsonReaderFromFile(path string) (*JsonReader, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	re := NewJsonReader(fi)
	re.closer = fi
	return re, nil
}

func NewJsonReaderFromString(rawData string) *JsonReader {
	re := strings.NewReader(rawData)
	return NewJsonReader(re)
}

func NewJsonReader(r io.Reader) *JsonReader {
	return &JsonReader{
		reader: r,
	}
}

func (r *JsonReader) Read(v any) error {
	if r.closer != nil {
		defer r.closer.Close()
	}

	decoder := json.NewDecoder(r.reader)
	for decoder.More() {
		err := decoder.Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}
