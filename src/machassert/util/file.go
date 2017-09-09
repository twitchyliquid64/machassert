package util

import (
	"os"
	"path"
)

// PathSanitize performs shell-like expansion for characters like ~.
func PathSanitize(in string) string {
	if len(in) > 0 && in[0] == '~' {
		return path.Join(os.Getenv("HOME"), in[2:])
	}
	return in
}
