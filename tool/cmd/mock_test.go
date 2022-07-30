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
	readContent     string
	writeBuffer     *bytes.Buffer
	isStdinEmpty    bool
	openFileError   error
	createFileError error
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		isStdinEmpty: true,
	}
}

func (r *mockRepository) GetFileReader(path string) (io.ReadCloser, error) {
	if r.openFileError != nil {
		return nil, r.openFileError
	}
	re := strings.NewReader(r.readContent)
	recl := io.NopCloser(re)
	return recl, nil
}

func (r *mockRepository) GetStdinReader() io.ReadCloser {
	re := strings.NewReader(r.readContent)
	recl := io.NopCloser(re)
	return recl
}

func (r *mockRepository) IsStdinEmpty() bool {
	return r.isStdinEmpty
}

func (r *mockRepository) CreateFileWriter(path string) (io.WriteCloser, error) {
	if r.createFileError != nil {
		return nil, r.createFileError
	}
	r.writeBuffer = &bytes.Buffer{}
	return &WriteNopCloser{Writer: r.writeBuffer}, nil
}
