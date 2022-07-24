package utils

import (
	"os"
)

func IsStdinEmpty() bool {
	fi := os.Stdin
	info, err := fi.Stat()
	if err != nil {
		return true
	}
	return info.Size() == 0
}
