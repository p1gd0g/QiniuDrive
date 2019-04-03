package gui

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/comm"
	"github.com/p1gd0g/ui"
)

// FileWindow creates a new file window.
func FileWindow(
	accessKey *ui.Entry,
	secretKey *ui.Entry,
	bucket *ui.Entry,
	domain *ui.Entry,
	zone *ui.Combobox,
	window *ui.Window,
	fileList *FileList) {

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

	fileBar := ui.NewProgressBar()
	fileBar.Hide()
	fileBar.SetValue(-1)

	fileVBox := ui.NewVerticalBox()
	fileVBox.SetPadded(true)
	fileVBox.Append(ui.NewLabel("文件信息"), false)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileList.HBox, false)
	fileVBox.Append(ui.NewHorizontalBox(), true)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileOpHBox, false)
	fileVBox.Append(fileBar, false)

	window.SetMargined(true)
	window.SetChild(fileVBox)

	fileUp.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Upload.")

		file := ui.OpenFile(window)
		if file == "" {
			return
		}
		log.Println("File is", file)

		fileBar.Show()
		log.Println("Bar shows.")

		go func() {

			err := comm.Upload(
				accessKey.Text(), secretKey.Text(), bucket.Text(),
				file, zone.Selected())

			ui.QueueMain(func() {
				if err != nil {
					ui.MsgBoxError(window, "Error!", err.Error())
					return
				}
				log.Println("Upload successfully.")

				err = fileList.Display(
					accessKey, secretKey, bucket)
				if err != nil {
					ui.MsgBoxError(window, "Error!", err.Error())
				}
				log.Println("Display successfully.")

				fileBar.Hide()
				log.Println("Bar hides.")
			})
		}()

	})

	fileDn.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Download.")

		fileBar.Show()
		log.Println("Bar shows.")

		go func() {
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
			fileBar.Hide()
			log.Println("Bar hides.")
		}()
	})

	fileDl.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Delete.")

		fileBar.Show()
		log.Println("Bar shows.")

		go func() {
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

			ui.QueueMain(func() {

				err := fileList.Display(
					accessKey, secretKey, bucket)
				if err != nil {
					ui.MsgBoxError(window, "Error!", err.Error())
				}
				fileBar.Hide()
				log.Println("Bar hides.")
			})
		}()
	})

	fileRd.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Remote download.")

		URLWindow(accessKey, secretKey, bucket,
			fileList, fileBar)

	})

	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
}
