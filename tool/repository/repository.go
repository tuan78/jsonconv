package repository

import (
	"io"
	"os"
)

type (
	Repository interface {
		GetFileReadCloser(path string) (io.ReadCloser, error)

		GetStdinReadCloser() io.ReadCloser

		IsStdinEmpty() bool
	}

	repository struct{}
)

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetFileReadCloser(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (r *repository) GetStdinReadCloser() io.ReadCloser {
	return os.Stdin
}

func (r *repository) IsStdinEmpty() bool {
	fi := os.Stdin
	info, err := fi.Stat()
	if err != nil {
		return true
	}
	return info.Size() == 0
}
