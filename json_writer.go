package jsonconv

import (
	"encoding/json"
	"io"
)

// A JsonWriter writes JSON values to an output stream.
type JsonWriter struct {
	// EscapeHTML specifies whether problematic HTML characters
	// should be escaped inside JSON quoted strings.
	// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
	// to avoid certain safety problems that can arise when embedding JSON in HTML.
	EscapeHTML bool

	writer io.Writer
}

// NewJsonWriter returns a new JsonWriter that writes to w.
func NewJsonWriter(w io.Writer) *JsonWriter {
	return &JsonWriter{
		writer:     w,
		EscapeHTML: true,
	}
}

// Write writes the JSON encoding of v to the stream,
// followed by a newline character.
func (r *JsonWriter) Write(v any) error {
	encoder := json.NewEncoder(r.writer)
	encoder.SetEscapeHTML(r.EscapeHTML)
	return encoder.Encode(v)
}
