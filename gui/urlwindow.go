package gui

import (
	"github.com/p1gd0g/ui"
)

// URLWindow creates a url window.
func URLWindow() (
	entry *ui.Entry, button *ui.Button, window *ui.Window) {
	entry = ui.NewEntry()

	form := ui.NewForm()
	form.SetPadded(true)
	form.Append("url", entry, true)

	button = ui.NewButton("确认")

	hBox := ui.NewHorizontalBox()
	hBox.SetPadded(true)
	hBox.Append(form, true)
	hBox.Append(button, false)

	window = ui.NewWindow("url", 1, 1, false)
	window.SetMargined(true)
	window.SetChild(hBox)

	window.OnClosing(func(*ui.Window) bool {
		return true
	})

	window.Show()
	return
}
