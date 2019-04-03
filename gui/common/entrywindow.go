package common

import (
	"github.com/p1gd0g/ui"
)

// NewEntryWindow creates a window with an entry.
func NewEntryWindow(s string) (
	*ui.Window, *ui.Entry, *ui.Button) {

	entry := ui.NewEntry()

	form := ui.NewForm()
	form.SetPadded(true)
	form.Append(s, entry, true)

	button := ui.NewButton("确认")

	hBox := ui.NewHorizontalBox()
	hBox.SetPadded(true)
	hBox.Append(form, true)
	hBox.Append(button, false)

	window := ui.NewWindow(s, 1, 1, false)
	window.SetMargined(true)
	window.SetChild(hBox)

	window.OnClosing(func(*ui.Window) bool {
		return true
	})

	window.Show()

	return window, entry, button
}
