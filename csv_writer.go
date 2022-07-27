package jsonconv

import (
	"encoding/csv"
	"io"
)

// A CsvWriter writes records using CSV encoding.
type CsvWriter struct {
	Delimiter *rune // Field delimiter. If nil, it uses default value from csv.NewWriter
	UseCRLF   bool  // True to use \r\n as the line terminator
	writer    io.Writer
}

// NewCsvWriter returns a new CsvWriter that writes to w.
func NewCsvWriter(w io.Writer) *CsvWriter {
	return &CsvWriter{
		writer: w,
	}
}

// NewDelimiter returns a pointer to v.
func NewDelimiter(v rune) *rune {
	return &v
}

// Write writes all CSV data to w.
func (w *CsvWriter) Write(data CsvData) error {
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
