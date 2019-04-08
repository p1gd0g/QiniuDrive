package gui

import (
	"flag"
	"log"
	"strconv"

	"github.com/p1gd0g/QiniuDrive/gui/common"
	"github.com/p1gd0g/ui"
)

// LoginWindow creates the new login window.
func LoginWindow() {
	var ak = flag.Lookup("ak").Value.String()
	var sk = flag.Lookup("sk").Value.String()
	var bk = flag.Lookup("bk").Value.String()
	var dm = flag.Lookup("dm").Value.String()
	var zn = flag.Lookup("zn").Value.String()

	fileList := NewFileList()

	accessKey := ui.NewEntry()
	accessKey.SetText(ak)
	secretKey := ui.NewPasswordEntry()
	secretKey.SetText(sk)
	bucket := ui.NewEntry()
	bucket.SetText(bk)
	domain := ui.NewEntry()
	domain.SetText(dm)

	zone := common.NewCombobox(
		"Huadong", "Huabei", "Huanan", "Beimei")
	zoneIndex, _ := strconv.Atoi(zn)
	zone.SetSelected(zoneIndex)

	loginForm := ui.NewForm()
	loginForm.SetPadded(true)
	loginForm.Append("accessKey", accessKey, false)
	loginForm.Append("secretKey", secretKey, false)
	loginForm.Append("bucket", bucket, false)
	loginForm.Append("domain", domain, false)
	loginForm.Append("zone", zone, false)

	loginGroup := ui.NewGroup("user info")
	loginGroup.SetMargined(true)
	loginGroup.SetChild(loginForm)

	loginButton := ui.NewButton("login")

	loginBar := ui.NewProgressBar()
	loginBar.Hide()
	loginBar.SetValue(-1)

	loginVBox := ui.NewVerticalBox()
	loginVBox.SetPadded(true)
	loginVBox.Append(loginGroup, false)
	loginVBox.Append(loginButton, false)
	loginVBox.Append(loginBar, false)

	loginWindow := ui.NewWindow("login", 200, 1, false)
	loginWindow.SetMargined(true)
	loginWindow.SetChild(loginVBox)

	loginButton.OnClicked(func(*ui.Button) {
		loginBar.Show()
		log.Println("accessKey:", accessKey.Text())
		log.Println("secretKey:", secretKey.Text())

		// go func() {
		err := fileList.Display(
			accessKey, secretKey, bucket)

		// ui.QueueMain(func() {
		if err != nil {
			ui.MsgBoxError(loginWindow, "Error!", err.Error())
			loginBar.Hide()
			return
		}
		log.Println("List files successfully.")

		loginBar.Hide()
		log.Println("loginBar hided.")

		loginWindow.Hide()
		log.Println("loginWindow hided.")

		FileWindow(accessKey, secretKey, bucket, domain,
			zone, fileList)
		// })
		// }()
	})

	bucket.OnChanged(func(*ui.Entry) {
		domain.SetText(bucket.Text())
	})

	loginWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		log.Println("loginWindow quitted.")

		return true
	})

	loginWindow.Show()
	log.Println("loginWindow showed.")
}
