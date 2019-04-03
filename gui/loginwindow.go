package gui

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/gui/common"
	"github.com/p1gd0g/ui"
)

// LoginWindow creates the new login window.
func LoginWindow() (
	*ui.Entry, *ui.Entry, *ui.Entry, *ui.Entry, *ui.Combobox,
	*ui.Window, *FileList) {
	fileList := NewFileList()

	accessKey := ui.NewEntry()
	secretKey := ui.NewPasswordEntry()
	bucket := ui.NewEntry()
	domain := ui.NewEntry()

	zone := common.NewCombobox("华东", "华北", "华南", "北美")
	zone.SetSelected(2)

	loginForm := ui.NewForm()
	loginForm.SetPadded(true)
	loginForm.Append("accessKey", accessKey, false)
	loginForm.Append("secretKey", secretKey, false)
	loginForm.Append("bucket", bucket, false)
	loginForm.Append("domain", domain, false)
	loginForm.Append("zone", zone, false)

	loginGroup := ui.NewGroup("登录信息")
	loginGroup.SetMargined(true)
	loginGroup.SetChild(loginForm)

	loginButton := ui.NewButton("登录")

	loginBar := ui.NewProgressBar()
	loginBar.Hide()
	loginBar.SetValue(-1)

	loginVBox := ui.NewVerticalBox()
	loginVBox.SetPadded(true)
	loginVBox.Append(loginGroup, false)
	loginVBox.Append(loginButton, false)
	loginVBox.Append(loginBar, false)

	login := ui.NewWindow("登录", 200, 1, false)
	login.SetMargined(true)
	login.SetChild(loginVBox)

	window := ui.NewWindow("QiniuDrive", 600, 600, false)

	loginButton.OnClicked(func(*ui.Button) {
		loginBar.Show()
		log.Println("accessKey:", accessKey.Text())
		log.Println("secretKey:", secretKey.Text())

		go func() {
			err := fileList.Display(
				accessKey.Text(), secretKey.Text(), bucket.Text())

			ui.QueueMain(func() {
				if err == nil {
					log.Println("List files successfully.")
					login.Hide()

					window.Show()
				} else {
					ui.MsgBoxError(login, "Error!", err.Error())
				}
				loginBar.Hide()
			})
		}()
	})

	bucket.OnChanged(func(*ui.Entry) {
		domain.SetText(bucket.Text())
	})

	login.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	login.Show()

	return accessKey, secretKey, bucket, domain,
		zone, window, fileList
}
