package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/p1gd0g/QiniuDrive/gui"
	"github.com/p1gd0g/QiniuDrive/tool"
	"github.com/p1gd0g/ui"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"

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

	fileNameVBox := ui.NewVerticalBox()
	fileMTypeVBox := ui.NewVerticalBox()
	fileSizeVBox := ui.NewVerticalBox()
	fileCheckboxVBox := ui.NewVerticalBox()

	fileHBox := ui.NewHorizontalBox()
	fileHBox.SetPadded(true)
	fileHBox.Append(fileNameVBox, true)
	fileHBox.Append(ui.NewHorizontalSeparator(), false)
	fileHBox.Append(fileMTypeVBox, false)
	fileHBox.Append(ui.NewHorizontalSeparator(), false)
	fileHBox.Append(fileSizeVBox, false)
	fileHBox.Append(ui.NewHorizontalSeparator(), false)
	fileHBox.Append(fileCheckboxVBox, false)

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
	fileVBox.Append(fileHBox, false)
	fileVBox.Append(ui.NewHorizontalBox(), true)
	fileVBox.Append(ui.NewVerticalSeparator(), false)
	fileVBox.Append(fileOpHBox, false)

	window := ui.NewWindow("QiniuDrive", 600, 600, false)
	window.SetMargined(true)
	window.SetChild(fileVBox)

	loginButton.OnClicked(func(*ui.Button) {
		log.Println("accessKey:", accessKey.Text())
		log.Println("secretKey:", secretKey.Text())

		mac := auth.New(accessKey.Text(), secretKey.Text())
		log.Println("mac created.")

		cfg := storage.Config{}

		bucketManager := storage.NewBucketManager(mac, &cfg)
		log.Println("bucketManager created.")

		log.Println("Listing files.")
		fileNameList, fileCheckboxList, err :=
			refresh(bucketManager, bucket.Text(),
				fileNameVBox, fileMTypeVBox, fileSizeVBox,
				fileCheckboxVBox)

		if err == nil {
			log.Println("List files successfully.")
			login.Hide()

			switch zone.Selected() {
			case 0:
				cfg.Zone = &storage.ZoneHuadong
				log.Println("Zone: Huadong.")
			case 1:
				cfg.Zone = &storage.ZoneHuabei
				log.Println("Zone: Huabei.")
			case 2:
				cfg.Zone = &storage.ZoneHuanan
				log.Println("Zone: Huanan.")
			case 3:
				cfg.Zone = &storage.ZoneBeimei
				log.Println("Zone: Beimei.")
			default:
				log.Println("No zone!")
			}

			fileUp.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Upload.")

				file := ui.OpenFile(window)
				fileName := tool.GetFileName(file)

				putPolicy := storage.PutPolicy{
					Scope: bucket.Text() + ":" + fileName,
				}
				upToken := putPolicy.UploadToken(mac)

				formUploader := storage.NewFormUploader(&cfg)
				ret := storage.PutRet{}

				err := formUploader.PutFile(context.Background(),
					&ret, upToken, fileName, file, &storage.PutExtra{})

				if err != nil {
					ui.MsgBoxError(window, "Error!", err.Error())
					return
				}

				log.Println("Upload successfully.")

				fileInfo, sErr := bucketManager.Stat(bucket.Text(),
					fileName)
				if sErr != nil {
					ui.MsgBoxError(window, "Error!", sErr.Error())
					return
				}

				log.Println("Fetch the uploaded file's info successfully.")

				fileNameList = append(fileNameList, fileName)
				log.Println("File name list increased.")

				fileNameVBox.Append(ui.NewLabel(fileName), true)
				log.Println("Added the file name.")

				fileMTypeVBox.Append(
					ui.NewLabel(fileInfo.MimeType), true)
				log.Println("Added the file mime type.")

				fileSizeVBox.Append(
					ui.NewLabel(formatSize(fileInfo.Fsize)), true)
				log.Println("Added the file size.")

				tempCheckbox := ui.NewCheckbox("")
				fileCheckboxList =
					append(fileCheckboxList, *tempCheckbox)
				log.Println("File checkbox list increased.")

				fileCheckboxVBox.Append(tempCheckbox, true)
				log.Println("Added the file checkbox.")
				log.Println("Added a new row.")

			})

			fileDn.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Download.")

				for index := 0; index < len(fileNameList); index++ {
					if fileCheckboxList[index].Checked() {
						go func(i int) {

							out, err := os.Create(fileNameList[i])
							defer out.Close()
							if err != nil {
								ui.MsgBoxError(window, "Error!",
									err.Error())
								return
							}

							resp, err := http.Get("http://" +
								domain.Text() + "/" + fileNameList[i])
							if err != nil {
								ui.MsgBoxError(window, "Error!",
									err.Error())
								return
							}
							defer resp.Body.Close()

							_, err = io.Copy(out, resp.Body)
							if err != nil {
								ui.MsgBoxError(window, "Error!",
									err.Error())
							} else {
								log.Println("Download",
									fileNameList[i], "successfully.")
							}

						}(index)
					}
				}
			})

			fileDl.OnClicked(func(*ui.Button) {
				log.Println("Button clicked: Delete.")

				for index := 0; index < len(fileNameList); index++ {
					if fileCheckboxList[index].Checked() {
						log.Println("To be deleted:", fileNameList[index])

						func(i int) {
							err = bucketManager.Delete(bucket.Text(),
								fileNameList[i])

							if err != nil {
								ui.MsgBoxError(window, "Error!", err.Error())
							} else {
								log.Println("Delete one file successfully.")
							}
						}(index)
					}
				}

				log.Println("All selected files deleted.")

				fileNameList, fileCheckboxList, err =
					refresh(bucketManager, bucket.Text(),
						fileNameVBox, fileMTypeVBox, fileSizeVBox,
						fileCheckboxVBox)

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
					fileName :=
						url.Text()[strings.LastIndex(url.Text(), "/")+1:]
					fetchRet, err :=
						bucketManager.Fetch(url.Text(),
							bucket.Text(), fileName)

					if err != nil {
						ui.MsgBoxError(window, "Error!", err.Error())
						return
					}

					log.Println("Remote download successfully.")

					fileNameVBox.Append(ui.NewLabel(fileName), true)
					fileMTypeVBox.Append(
						ui.NewLabel(fetchRet.MimeType), true)
					fileSizeVBox.Append(
						ui.NewLabel(formatSize(fetchRet.Fsize)), true)

					fileNameList = append(fileNameList, fetchRet.Key)
					tempCheckbox := ui.NewCheckbox("")
					fileCheckboxList =
						append(fileCheckboxList, *tempCheckbox)
					fileCheckboxVBox.Append(tempCheckbox, true)

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

func refresh(
	bucketManager *storage.BucketManager,
	bucket string,
	fileNameVBox *ui.Box,
	fileMTypeVBox *ui.Box,
	fileSizeVBox *ui.Box,
	fileCheckboxVBox *ui.Box,
) (
	fileNameList []string,
	fileCheckboxList []ui.Checkbox,
	err error,
) {
	log.Println("Refreshing the file list.")

	fileNameVBox.Clear()
	fileMTypeVBox.Clear()
	fileSizeVBox.Clear()
	fileCheckboxVBox.Clear()

	marker := ""
	for {

		entries, _, nextMarker, hashNext, sErr :=
			bucketManager.ListFiles(bucket,
				"", "", marker, 1000)

		if sErr != nil {
			err = sErr
			break
		} else {
			log.Println("Fetch file info successfully.")
		}

		for _, entry := range entries {

			fileNameList = append(fileNameList, entry.Key)
			fileNameVBox.Append(ui.NewLabel(entry.Key), true)
			fileMTypeVBox.Append(ui.NewLabel(entry.MimeType), true)
			fileSizeVBox.Append(
				ui.NewLabel(formatSize(entry.Fsize)), true)
			tempCheckbox := ui.NewCheckbox("")
			fileCheckboxList =
				append(fileCheckboxList, *tempCheckbox)
			fileCheckboxVBox.Append(tempCheckbox, true)
		}

		if hashNext {
			marker = nextMarker
		} else {
			break
		}

	}
	log.Println("Refreshed the file list.")
	return
}

func formatSize(n int64) string {
	if n < 1024 {
		return strconv.FormatInt(n, 10) + " B"
	}

	nf := float64(n)

	nf /= 1024
	if nf < 1024 {
		return strconv.FormatFloat(nf, 'f', 2, 64) + " KB"
	}

	nf /= 1024
	if nf < 1024 {
		return strconv.FormatFloat(nf, 'f', 2, 64) + " MB"
	}

	nf /= 1024
	if nf < 1024 {
		return strconv.FormatFloat(nf, 'f', 2, 64) + " GB"
	}

	nf /= 1024
	return strconv.FormatFloat(nf, 'f', 2, 64) + " TB"
}
