// +build !windows

package gui

import (
	"github.com/p1gd0g/ui"
)

// Start starts the gui.
func Start() {
	err := ui.Main(func() {

		LoginWindow()

	})
	if err != nil {
		panic(err)
	}
}
