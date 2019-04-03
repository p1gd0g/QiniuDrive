// +build windows

package gui

import (
	"github.com/p1gd0g/ui"

	// Only for Windows.
	_ "github.com/andlabs/ui/winmanifest"
)

// Start starts the ui.
func Start() {
	err := ui.Main(func() {
		accessKey, secretKey, bucket, domain,
			zone, window, fileList :=
			LoginWindow()

		FileWindow(accessKey, secretKey, bucket, domain,
			zone, window, fileList)
	})
	if err != nil {
		panic(err)
	}
}
