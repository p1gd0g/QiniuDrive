// +build windows

package tool

import (
	"strings"
)

// GetFileName gets the name of a fullpath.
func GetFileName(file string) (name string) {

	name = file[strings.LastIndex(file, "\\")+1:]

	return
}
