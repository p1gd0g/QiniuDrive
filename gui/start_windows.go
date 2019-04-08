// +build windows

package gui

import (
	"github.com/p1gd0g/ui"

	// Only for Windows.
	_ "github.com/andlabs/ui/winmanifest"
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
