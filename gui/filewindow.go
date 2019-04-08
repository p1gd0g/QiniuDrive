package gui

import (
	"log"

	"github.com/p1gd0g/QiniuDrive/comm"
	"github.com/p1gd0g/QiniuDrive/crypto"
	"github.com/p1gd0g/ui"
)

// FileWindow creates a new file window.
func FileWindow(
	accessKey *ui.Entry,
	secretKey *ui.Entry,
	bucket *ui.Entry,
	domain *ui.Entry,
	zone *ui.Combobox,
	fileList *FileList) {

	fileUp := ui.NewButton("Upload")
	fileDn := ui.NewButton("Download")
	fileDl := ui.NewButton("Delete")
	fileFc := ui.NewButton("Fetch")

	fileOpHBox := ui.NewHorizontalBox()
	fileOpHBox.SetPadded(true)
	fileOpHBox.Append(fileUp, true)
	fileOpHBox.Append(fileDn, true)
	fileOpHBox.Append(fileDl, true)
	fileOpHBox.Append(fileFc, true)

	enc := ui.NewButton("Encrypt")
	dec := ui.NewButton("Decrypt")

	keyEntry := ui.NewEntry()

	cryptoHBox := ui.NewHorizontalBox()
	cryptoHBox.SetPadded(true)
	cryptoHBox.Append(enc, true)
	cryptoHBox.Append(ui.NewLabel("or"), false)
	cryptoHBox.Append(dec, true)
	cryptoHBox.Append(ui.NewLabel("a file with the key:"), false)
	cryptoHBox.Append(keyEntry, true)

	fileBar := ui.NewProgressBar()
	fileBar.Hide()
	fileBar.SetValue(-1)

	fileVBox := ui.NewVerticalBox()
	fileVBox.SetPadded(true)
	fileVBox.Append(ui.NewLabel("File info"), false)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileList.HBox, false)
	fileVBox.Append(ui.NewHorizontalBox(), true)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileOpHBox, false)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(cryptoHBox, false)
	fileVBox.Append(fileBar, false)

	fileWindow := ui.NewWindow("Qiniu Drive v0.1", 600, 600, false)
	fileWindow.SetMargined(true)
	fileWindow.SetChild(fileVBox)

	fileUp.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Upload.")

		file := ui.OpenFile(fileWindow)
		if file == "" {
			log.Println("Empty file.")

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

	fileFc.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Fetch.")

		urlEntry, urlButton, URLWindow := URLWindow()

		urlButton.OnClicked(func(*ui.Button) {
			URLWindow.Hide()
			fileBar.Show()
			log.Println("Bar shows.")

			go func() {
				err := comm.Fetch(
					accessKey.Text(),
					secretKey.Text(),
					bucket.Text(),
					urlEntry.Text())

				if err != nil {
					return
				}
				log.Println("Fetch successfully.")

				ui.QueueMain(func() {

					if err != nil {
						ui.MsgBoxError(fileWindow, "Error!", err.Error())
					}
					err = fileList.Display(
						accessKey, secretKey, bucket)
					if err != nil {
						ui.MsgBoxError(fileWindow, "Error!", err.Error())
					}
					fileBar.Hide()
					log.Println("Bar hides.")
				})
			}()
		})
	})

	enc.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Enc.")

		file := ui.OpenFile(fileWindow)
		if file == "" || keyEntry.Text() == "" {
			log.Println("Empty file or key.")

			return
		}
		log.Println("File is", file)

		fileBar.Show()
		log.Println("Bar shows.")

		go func() {
			err := crypto.Encrypt(keyEntry.Text(), file)
			if err != nil {
				ui.MsgBoxError(fileWindow, "Error!", err.Error())
			}
			log.Println("Encrypt successfully.")

			fileBar.Hide()
			log.Println("Bar hides.")
		}()
	})

	dec.OnClicked(func(*ui.Button) {
		log.Println("Button clicked: Dec.")

		file := ui.OpenFile(fileWindow)
		if file == "" || keyEntry.Text() == "" {
			log.Println("Empty file or key.")

			return
		}
		log.Println("File is", file)

		fileBar.Show()
		log.Println("Bar shows.")

		go func() {
			err := crypto.Decrypt(keyEntry.Text(), file)
			if err != nil {
				ui.MsgBoxError(fileWindow, "Error!", err.Error())
			}
			log.Println("Decrypt successfully.")

			fileBar.Hide()
			log.Println("Bar hides.")
		}()
	})

	fileWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	fileWindow.Show()
	log.Println("fileWindow showed.")
}
