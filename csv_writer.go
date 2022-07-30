package jsonconv

import (
	"encoding/csv"
	"io"
)

// A CsvWriter writes records using CSV encoding.
type CsvWriter struct {
	// Field delimiter. Set to ',' (CsvComma) by default in NewCsvWriter
	Delimiter rune

	// True to use \r\n as the line terminator
	UseCRLF bool

	writer io.Writer
}

const CsvComma rune = ','

// NewCsvWriter returns a new CsvWriter that writes to w.
func NewCsvWriter(w io.Writer) *CsvWriter {
	return &CsvWriter{
		writer:    w,
		Delimiter: CsvComma,
	}
}

// Write writes all CSV data to w.
func (w *CsvWriter) Write(data CsvData) error {
	writer := csv.NewWriter(w.writer)
	writer.Comma = w.Delimiter
	writer.UseCRLF = w.UseCRLF

	defer writer.Flush()
	for _, v := range data {
		if err := writer.Write(v); err != nil {
			return err
		}
	}
	return nil
}
