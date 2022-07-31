package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// WriteNopCloser .
type WriteNopCloser struct {
	io.Writer
}

func (WriteNopCloser) Close() error {
	return nil
}

// Mock Logger.
type mockLogger struct {
	msg string
}

func NewMockLogger() *mockLogger {
	return &mockLogger{}
}

func (l *mockLogger) Printf(format string, i ...interface{}) {
	l.msg = fmt.Sprintf(format, i...)
}

// Mock Repository.
type mockRepository struct {
	readerContent     string
	writerBuffer      *bytes.Buffer
	isStdinEmpty      bool
	fileOpeningError  error
	fileCreatingError error
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		isStdinEmpty: true,
	}
}

func (r *mockRepository) GetFileReader(path string) (io.ReadCloser, error) {
	if r.fileOpeningError != nil {
		return nil, r.fileOpeningError
	}
	re := strings.NewReader(r.readerContent)
	recl := io.NopCloser(re)
	return recl, nil
}

func (r *mockRepository) GetStdinReader() io.ReadCloser {
	re := strings.NewReader(r.readerContent)
	recl := io.NopCloser(re)
	return recl
}

func (r *mockRepository) IsStdinEmpty() bool {
	return r.isStdinEmpty
}

func (r *mockRepository) CreateFileWriter(path string) (io.WriteCloser, error) {
	if r.fileCreatingError != nil {
		return nil, r.fileCreatingError
	}
	r.writerBuffer = &bytes.Buffer{}
	return &WriteNopCloser{Writer: r.writerBuffer}, nil
}
