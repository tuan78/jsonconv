package utils

import (
	"path/filepath"
)

// IsFilePath checks whether or not v is file path.
// For example:
//	- usr//sample.csv
//	- ~/file.json
//  - /dir/contains/file
func IsFilePath(v string) bool {
	dir := filepath.Dir(v)
	return dir != "."
}
