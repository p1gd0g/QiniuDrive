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
	fileWindow *ui.Window,
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

	fileWindow.SetChild(fileVBox)

	fileUp.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Upload.")

		file := ui.OpenFile(fileWindow)
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
					ui.MsgBoxError(fileWindow, "Error!", err.Error())
					return
				}
				log.Println("Upload successfully.")

				err = fileList.Display(
					accessKey, secretKey, bucket)
				if err != nil {
					ui.MsgBoxError(fileWindow, "Error!", err.Error())
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
						ui.QueueMain(func() {
							if err != nil {
								ui.MsgBoxError(fileWindow, "Error!",
									err.Error())
							}

							fileBar.Hide()
							log.Println("Bar hides.")
						})
					}(fileList.NameList[i], domain.Text())
				}
			}

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

					ui.QueueMain(func() {
						if err != nil {
							ui.MsgBoxError(fileWindow, "Error!", err.Error())
							return
						}
						log.Println("Delete one file successfully.")

						err := fileList.Display(
							accessKey, secretKey, bucket)
						if err != nil {
							ui.MsgBoxError(fileWindow, "Error!", err.Error())
						}
						fileBar.Hide()
						log.Println("Bar hides.")
					})
				}
			}
			log.Println("All selected files deleted.")

		}()
	})

	fileRd.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Remote download.")

		fileBar.Show()
		log.Println("Bar shows.")

		go func() {

			var err *error
			URLWindow(accessKey, secretKey, bucket,
				fileList, err)

			ui.QueueMain(func() {

				if err != nil {
					ui.MsgBoxError(fileWindow, "Error!", (*err).Error())
				}
				*err = fileList.Display(
					accessKey, secretKey, bucket)
				if err != nil {
					ui.MsgBoxError(fileWindow, "Error!", (*err).Error())
				}
				fileBar.Hide()
				log.Println("Bar hides.")
			})
		}()

	})

	fileWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
}
