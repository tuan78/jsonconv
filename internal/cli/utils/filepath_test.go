package utils

import (
	"path/filepath"
	"testing"
)

func TestIsFilePath(t *testing.T) {
	// Test path should be a file path.
	paths := []string{
		"usr" + string(filepath.Separator) + string(filepath.Separator) + "sample.csv",
		filepath.Join("~", "file.json"),
		filepath.Join("dir", "contains", "file"),
		filepath.Join("dir", "file.json"),
		filepath.Join("dir", "file.x"),
	}
	for _, v := range paths {
		if !IsFilePath(v) {
			t.Fatalf("%s should be a file path", v)
		}
	}

	// Test path should not be a file path.
	paths = []string{
		filepath.Join(".", "file.json"),
		"file",
		"file.json",
		"file.x",
	}
	for _, v := range paths {
		if IsFilePath(v) {
			t.Fatalf("%s should not be a file path", v)
		}
	}
}
