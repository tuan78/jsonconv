package jsonconv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
)

type CsvWriter struct {
	Delimiter *rune // Field delimiter. If nil, it uses default value from csv.NewWriter
	UseCRLF   bool  // True to use \r\n as the line terminator
	writer    io.Writer
	closer    io.Closer
}

func NewCsvWriterFromFile(path string) (*CsvWriter, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}
	fi, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	wr := NewCsvWriter(fi)
	wr.closer = fi
	return NewCsvWriter(fi), nil
}

func NewCsvWriterFromByteBuffer() (*CsvWriter, *bytes.Buffer) {
	buf := bytes.NewBuffer([]byte{})
	wr := NewCsvWriter(buf)
	return wr, buf
}

func NewCsvWriter(w io.Writer) *CsvWriter {
	return &CsvWriter{
		writer: w,
	}
}

func (w *CsvWriter) Write(data CsvData) error {
	if w.closer != nil {
		defer w.closer.Close()
	}

	writer := csv.NewWriter(w.writer)
	if w.Delimiter != nil {
		writer.Comma = *w.Delimiter
	}
	writer.UseCRLF = w.UseCRLF

	defer writer.Flush()
	for _, v := range data {
		if err := writer.Write(v); err != nil {
			return err
		}
	}
	return nil
}
