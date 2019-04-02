package gui

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/comm"
	"github.com/p1gd0g/ui"
)

// Login creates the login window and inputs info.
func Login() {
	window := ui.NewWindow("QiniuDrive", 600, 600, false)
	fileList := NewFileList()

	accessKey, secretKey, bucket, domain, zone :=
		NewLoginWindow(window, fileList)

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

	window.SetMargined(true)
	window.SetChild(fileVBox)

	fileUp.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Upload.")

		file := ui.OpenFile(window)

		err := comm.Upload(
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
					err := comm.Download(name, domain)
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

				err := comm.Delete(
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

		err := fileList.Display(
			accessKey.Text(), secretKey.Text(), bucket.Text())
		if err != nil {
			ui.MsgBoxError(window, "Error!", err.Error())
		}
	})

	fileRd.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Remote download.")

		urlWindow, urlEntry, urlButton :=
			NewEntryWindow("url")

		urlButton.OnClicked(func(*ui.Button) {
			urlWindow.Hide()

			err :=
				comm.RemoteDownload(
					accessKey.Text(),
					secretKey.Text(),
					bucket.Text(),
					urlEntry.Text())

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

		})
		urlWindow.OnClosing(func(*ui.Window) bool {
			return true
		})
		urlWindow.Show()
	})

	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
}
