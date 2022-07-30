package cmd

import (
	"fmt"
	"io"
	"strings"
)

type mockCmdLogger struct {
	msg string
}

func NewMockCmdLogger() *mockCmdLogger {
	return &mockCmdLogger{}
}

func (l *mockCmdLogger) Printf(format string, i ...interface{}) {
	l.msg = fmt.Sprintf(format, i...)
}

type mockRepository struct {
	content      string
	isStdinEmpty bool
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		isStdinEmpty: true,
	}
}

func (r *mockRepository) GetFileReadCloser(path string) (io.ReadCloser, error) {
	re := strings.NewReader(r.content)
	recl := io.NopCloser(re)
	return recl, nil
}

func (r *mockRepository) GetStdinReadCloser() io.ReadCloser {
	re := strings.NewReader(r.content)
	recl := io.NopCloser(re)
	return recl
}

func (r *mockRepository) IsStdinEmpty() bool {
	return r.isStdinEmpty
}
