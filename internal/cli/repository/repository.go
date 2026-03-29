package repository

import (
	"io"
	"os"
	"path/filepath"

	"github.com/tuan78/jsonconv/v2/internal/cli/utils"
)

type (
	// A Repository used to interact with file system, database, networking and more.
	Repository interface {
		GetFileReader(path string) (io.ReadCloser, error)

		GetStdinReader() io.ReadCloser

		IsStdinEmpty() bool

		CreateFileWriter(path string) (io.WriteCloser, error)
	}

	repository struct{}
)

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetFileReader(path string) (io.ReadCloser, error) {
	// #nosec G304 -- CLI tool explicitly opens user-specified files
	return os.Open(filepath.Clean(path))
}

func (r *repository) GetStdinReader() io.ReadCloser {
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

func (r *repository) CreateFileWriter(path string) (io.WriteCloser, error) {
	// Clean path to prevent traversal attacks
	path = filepath.Clean(path)

	// Check file path and make dir accordingly.
	if utils.IsFilePath(path) {
		// Ensure all dir in path exists with secure permissions.
		err := os.MkdirAll(filepath.Dir(path), 0750)
		if err != nil {
			return nil, err
		}
		// #nosec G304 -- CLI tool explicitly creates user-specified files
		fi, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		return fi, nil
	}

	// Path is only file name so override it with full path (working dir + file name).
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(dir, path)
	// #nosec G304 -- CLI tool explicitly creates user-specified files
	fi, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return fi, nil
}
