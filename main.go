package main

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/comm"
	"github.com/p1gd0g/QiniuDrive/gui"
	"github.com/p1gd0g/ui"

	_ "github.com/andlabs/ui/winmanifest"
)

func main() {
	log.SetFlags(log.Lshortfile)

	err := ui.Main(login)
	if err != nil {
		panic(err)
	}
}

func login() {

	accessKey := ui.NewEntry()
	secretKey := ui.NewPasswordEntry()
	bucket := ui.NewEntry()
	domain := ui.NewEntry()

	zone := gui.NewCombobox("华东", "华北", "华南", "北美")

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

	loginVBox := ui.NewVerticalBox()
	loginVBox.SetPadded(true)
	loginVBox.Append(loginGroup, false)
	loginVBox.Append(loginButton, false)

	login := ui.NewWindow("登录", 200, 1, false)
	login.SetMargined(true)
	login.SetChild(loginVBox)

	fileList := gui.NewFileList()

	fileUp := ui.NewButton("上传文件")
	fileDn := ui.NewButton("下载文件")
	fileDl := ui.NewButton("删除文件")
	fileRd := ui.NewButton("离线下载")

	fileOpHBox := ui.NewHorizontalBox()
	fileOpHBox.SetPadded(true)
	fileOpHBox.Append(fileUp, true)
	fileOpHBox.Append(fileDn, true)
	fileOpHBox.Append(fileDl, true)
	fileOpHBox.Append(fileRd, true)

	fileVBox := ui.NewVerticalBox()
	fileVBox.SetPadded(true)
	fileVBox.Append(ui.NewLabel("文件信息"), false)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileList.HBox, false)
	fileVBox.Append(ui.NewHorizontalBox(), true)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileOpHBox, false)

	window := ui.NewWindow("QiniuDrive", 600, 600, false)
	window.SetMargined(true)
	window.SetChild(fileVBox)

	loginButton.OnClicked(func(*ui.Button) {
		log.Println("accessKey:", accessKey.Text())
		log.Println("secretKey:", secretKey.Text())

		err := fileList.Display(
			accessKey.Text(), secretKey.Text(), bucket.Text())

		if err == nil {
			log.Println("List files successfully.")
			login.Hide()

			fileUp.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Upload.")

				file := ui.OpenFile(window)

				err = comm.Upload(
					accessKey.Text(), secretKey.Text(), bucket.Text(),
					file, zone.Selected())
				if err != nil {
					ui.MsgBoxError(window, "Error!", err.Error())
					return
				}

				log.Println("Upload successfully.")

				err = fileList.Display(
					accessKey.Text(), secretKey.Text(), bucket.Text())
				if err != nil {
					ui.MsgBoxError(window, "Error!",
						err.Error())
				}

			})

			fileDn.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Download.")

				for i := 0; i < len(fileList.NameList); i++ {
					if fileList.CheckboxList[i].Checked() {
						go func(name, domain string) {
							err = comm.Download(name, domain)
							if err != nil {
								ui.MsgBoxError(window, "Error!",
									err.Error())
							}
						}(fileList.NameList[i], domain.Text())
					}
				}
			})

			fileDl.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Delete.")

				for i := 0; i < len(fileList.NameList); i++ {
					if fileList.CheckboxList[i].Checked() {
						log.Println("To be deleted:", fileList.NameList[i])

						err = comm.Delete(
							accessKey.Text(),
							secretKey.Text(),
							bucket.Text(),
							fileList.NameList[i])

						if err != nil {
							ui.MsgBoxError(window, "Error!", err.Error())
						} else {
							log.Println("Delete one file successfully.")
						}

					}
				}
				log.Println("All selected files deleted.")

				err = fileList.Display(
					accessKey.Text(), secretKey.Text(), bucket.Text())
				if err != nil {
					ui.MsgBoxError(window, "Error!", err.Error())
				}
			})

			fileRd.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Remote download.")

				url := ui.NewEntry()
				urlButton := ui.NewButton("确定")

				urlHBox := ui.NewHorizontalBox()
				urlHBox.SetPadded(true)
				urlHBox.Append(ui.NewLabel("url"), false)
				urlHBox.Append(url, true)
				urlHBox.Append(urlButton, false)

				urlWindow := ui.NewWindow("url", 1, 1, false)
				urlWindow.SetMargined(true)
				urlWindow.SetChild(urlHBox)

				urlButton.OnClicked(func(*ui.Button) {
					urlWindow.Hide()

					go func() {
						err =
							comm.RemoteDownload(
								accessKey.Text(),
								secretKey.Text(),
								bucket.Text(),
								url.Text())

						if err != nil {
							ui.MsgBoxError(window, "Error!", err.Error())
							return
						}
						log.Println("Remote download successfully.")

						err = fileList.Display(
							accessKey.Text(), secretKey.Text(), bucket.Text())
						if err != nil {
							ui.MsgBoxError(window, "Error!", err.Error())
						}
					}()
				})
				urlWindow.OnClosing(func(*ui.Window) bool {
					return true
				})
				urlWindow.Show()
			})

			window.Show()
		} else {
			ui.MsgBoxError(login, "Error!", err.Error())
		}
	})

	login.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	login.Show()
}
