package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/andlabs/ui"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"

	_ "github.com/andlabs/ui/winmanifest"
)

func main() {
	err := ui.Main(func() {

		accessKey := ui.NewEntry()
		secretKey := ui.NewPasswordEntry()
		bucket := ui.NewEntry()
		domain := ui.NewEntry()

		zone := ui.NewCombobox()
		zone.Append("华东")
		zone.Append("华北")
		zone.Append("华南")
		zone.Append("北美")

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
		loginVBox.Append(loginGroup, false)
		loginVBox.Append(loginButton, false)

		login := ui.NewWindow("登录", 200, 1, false)
		login.SetMargined(true)
		login.SetChild(loginVBox)

		fileNameVBox := ui.NewVerticalBox()
		fileSizeVBox := ui.NewVerticalBox()
		fileSelectedVBox := ui.NewVerticalBox()

		fileHBox := ui.NewHorizontalBox()
		fileHBox.SetPadded(true)
		fileHBox.Append(fileNameVBox, false)
		fileHBox.Append(ui.NewHorizontalSeparator(), false)
		fileHBox.Append(fileSizeVBox, false)
		fileHBox.Append(ui.NewHorizontalSeparator(), false)
		fileHBox.Append(fileSelectedVBox, false)

		fileUp := ui.NewButton("上传文件")
		fileDn := ui.NewButton("下载文件")
		fileDl := ui.NewButton("删除文件")

		fileOpHBox := ui.NewHorizontalBox()
		fileOpHBox.SetPadded(true)
		fileOpHBox.Append(fileUp, false)
		fileOpHBox.Append(fileDn, false)
		fileOpHBox.Append(fileDl, false)

		fileVBox := ui.NewVerticalBox()
		fileVBox.SetPadded(true)
		fileVBox.Append(ui.NewLabel("文件信息"), false)
		fileVBox.Append(ui.NewVerticalSeparator(), false)
		fileVBox.Append(fileHBox, false)
		fileVBox.Append(ui.NewVerticalSeparator(), false)
		fileVBox.Append(fileOpHBox, false)

		window := ui.NewWindow("QiniuDrive", 400, 1, false)
		window.SetMargined(true)
		window.SetChild(fileVBox)

		loginButton.OnClicked(func(*ui.Button) {

			mac := auth.New(accessKey.Text(), secretKey.Text())

			cfg := storage.Config{}

			bucketManager := storage.NewBucketManager(mac, &cfg)
			marker := ""

			var loginError error
			var fileNameList []string
			var fileSelectedList []ui.Checkbox

			for {

				entries, _, nextMarker, hashNext, err :=
					bucketManager.ListFiles(bucket.Text(),
						"", "", marker, 1000)

				if err != nil {
					loginError = err
					break
				}

				for _, entry := range entries {

					fileNameList = append(fileNameList, entry.Key)
					fileNameVBox.Append(ui.NewLabel(entry.Key), true)
					fileSizeVBox.Append(
						ui.NewLabel(
							strconv.FormatInt(
								entry.Fsize, 10)), true)
					tempCheckbox := ui.NewCheckbox("")
					fileSelectedList =
						append(fileSelectedList, *tempCheckbox)
					fileSelectedVBox.Append(tempCheckbox, true)
				}

				if hashNext {
					marker = nextMarker
				} else {
					break
				}
			}

			if loginError == nil {
				login.Hide()

				switch zone.Selected() {
				case 0:
					cfg.Zone = &storage.ZoneHuadong
				case 1:
					cfg.Zone = &storage.ZoneHuabei
				case 2:
					cfg.Zone = &storage.ZoneHuanan
				case 3:
					cfg.Zone = &storage.ZoneBeimei
				}

				fileUp.OnClicked(func(*ui.Button) {

					file := ui.OpenFile(ui.NewWindow("选择文件",
						300, 300, false))
					fileName := file[strings.LastIndex(file, "/")+1:]

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
					} else {
						fileInfo, sErr := bucketManager.Stat(bucket.Text(),
							fileName)
						if sErr != nil {
							ui.MsgBoxError(window, "Error!", sErr.Error())
						} else {
							fileNameList = append(fileNameList, fileName)
							fileNameVBox.Append(ui.NewLabel(fileName), true)
							fileSizeVBox.Append(
								ui.NewLabel(
									strconv.FormatInt(
										fileInfo.Fsize, 10)), true)

							tempCheckbox := new(ui.Checkbox)
							fileSelectedList =
								append(fileSelectedList, *tempCheckbox)
							fileSelectedVBox.Append(tempCheckbox, true)
						}

					}

				})

				fileDn.OnClicked(func(*ui.Button) {

					for index := 0; index < len(fileNameList); index++ {
						if fileSelectedList[index].Checked() {
							go func(i int) {

								out, err := os.Create(fileNameList[i])
								defer out.Close()
								if err != nil {
									ui.MsgBoxError(login, "Error!",
										err.Error())
								} else {
									resp, err := http.Get("http://" +
										domain.Text() + "/" + fileNameList[i])
									if err != nil {
										ui.MsgBoxError(login, "Error!",
											err.Error())
									} else {
										_, err := io.Copy(out, resp.Body)
										if err != nil {
											ui.MsgBoxError(login, "Error!",
												err.Error())
										}
									}
									defer resp.Body.Close()
								}
							}(index)
						}
					}
				})

				window.Show()
			} else {
				ui.MsgBoxError(login, "Error!", loginError.Error())
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
	})
	if err != nil {
		panic(err)
	}
}
