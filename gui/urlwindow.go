package gui

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/comm"
	"github.com/p1gd0g/ui"
)

// URLWindow creates a url window.
func URLWindow(
	accessKey *ui.Entry,
	secretKey *ui.Entry,
	bucket *ui.Entry,
	fileList *FileList,
	err *error) {
	entry := ui.NewEntry()

	form := ui.NewForm()
	form.SetPadded(true)
	form.Append("url", entry, true)

	button := ui.NewButton("确认")

	hBox := ui.NewHorizontalBox()
	hBox.SetPadded(true)
	hBox.Append(form, true)
	hBox.Append(button, false)

	window := ui.NewWindow("url", 1, 1, false)
	window.SetMargined(true)
	window.SetChild(hBox)

	button.OnClicked(func(*ui.Button) {

		go func() {

			*err = comm.RemoteDownload(
				accessKey.Text(),
				secretKey.Text(),
				bucket.Text(),
				entry.Text())

			if err != nil {
				return
			}
			log.Println("Remote download successfully.")

			window.Hide()
		}()
	})

	window.OnClosing(func(*ui.Window) bool {
		return true
	})

	window.Show()
}
