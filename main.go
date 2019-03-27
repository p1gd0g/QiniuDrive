package main

import (
	"fmt"

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

		window := ui.NewWindow("QiniuDrive", 400, 200, false)

		loginButton.OnClicked(func(*ui.Button) {

			putPolicy := storage.PutPolicy{}

			mac := auth.New(accessKey.Text(), secretKey.Text())
			upToken := putPolicy.UploadToken(mac)
			fmt.Println(upToken)

			cfg := storage.Config{}

			bucketManager := storage.NewBucketManager(mac, &cfg)
			marker := ""

			var loginError error

			for {

				entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(bucket.Text(), "", "", marker, 1000)

				if err != nil {
					loginError = err
					break
				}

				for _, entry := range entries {
					fmt.Println(entry.Key)
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

				// client := http.Client{}

				window.Show()
			} else {
				ui.MsgBoxError(login, "Warning!", "Wrong user information!")
			}
		})

		// putPolicy := storage.PutPolicy{}
		// mac := auth.New(accessKey.Text(), secretKey.Text())
		// upToken := putPolicy.UploadToken(mac)

		// cfg := storage.Config{}

		// var path string
		// tab := ui.NewTab()

		// // tab1
		// buttonFile1 := ui.NewButton("选择文件")
		// box1 := ui.NewVerticalBox()
		// label1 := ui.NewLabel("未选择文件")
		// label2 := ui.NewLabel("未加密")
		// key1 := ui.NewEntry()
		// key1.SetText("输入密钥")
		// buttonEnc := ui.NewButton("加密")
		// box1.Append(buttonFile1, true)
		// box1.Append(label1, false)
		// box1.Append(key1, false)
		// box1.Append(buttonEnc, true)
		// box1.Append(label2, false)
		// tab.Append("加密", box1)

		// // tab2
		// box2 := ui.NewVerticalBox()
		// buttonFile2 := ui.NewButton("选择文件")
		// label3 := ui.NewLabel("未选择文件")
		// label4 := ui.NewLabel("未解密")
		// buttonDec := ui.NewButton("解密")
		// key2 := ui.NewEntry()
		// key2.SetText("输入密钥")
		// box2.Append(buttonFile2, true)
		// box2.Append(label3, false)
		// box2.Append(key2, false)
		// box2.Append(buttonDec, true)
		// box2.Append(label4, false)
		// tab.Append("解密", box2)

		// // tab3
		// buttonFile3 := ui.NewButton("选择文件")
		// box3 := ui.NewVerticalBox()
		// label5 := ui.NewLabel("未选择文件")
		// label6 := ui.NewLabel("未上传")
		// buttonUp := ui.NewButton("上传")
		// box3.Append(buttonFile3, true)
		// box3.Append(label5, false)
		// box3.Append(buttonUp, true)
		// box3.Append(label6, false)
		// tab.Append("上传", box3)

		// // tab4
		// box4 := ui.NewVerticalBox()
		// label7 := ui.NewLabel("未下载")
		// tab.Append("下载", box4)
		// list := file.List()
		// temp := ui.NewButton("")
		// for _, item := range list {
		// 	temp := ui.NewButton(item.Key)
		// 	temp.OnClicked(func(*ui.Button) {
		// 		fmt.Println(temp.Text())
		// 		file.Dn(temp.Text())
		// 		label7.SetText("已下载")
		// 	})
		// 	box4.Append(temp, false)
		// }
		// box4.Append(label7, false)

		// fileChooser := ui.NewWindow("选择文件", 500, 500, false)
		// window.SetChild(tab)

		// buttonFile1.OnClicked(func(*ui.Button) {
		// 	path = ui.OpenFile(filechooser)
		// 	label1.SetText("文件地址:" + path)

		// })
		// buttonEnc.OnClicked(func(*ui.Button) {
		// 	file, _ := os.Open(path)

		// 	crypto_p1gd0g.Enc(file, key1.Text())
		// 	label2.SetText("已加密")
		// })
		// buttonFile2.OnClicked(func(*ui.Button) {
		// 	path = ui.OpenFile(filechooser)
		// 	label3.SetText("文件地址:" + path)

		// })
		// buttonDec.OnClicked(func(*ui.Button) {
		// 	file, _ := os.Open(path)

		// 	crypto_p1gd0g.Dec(file, key2.Text())
		// 	label4.SetText("已解密")
		// })
		// buttonFile3.OnClicked(func(*ui.Button) {
		// 	path = ui.OpenFile(filechooser)
		// 	label5.SetText("文件地址:" + path)

		// })
		// buttonUp.OnClicked(func(*ui.Button) {
		// 	file.Up(path)
		// 	var j int
		// 	for i := 0; i < len(path); i++ {
		// 		if path[i] == '/' {
		// 			j = i
		// 		}
		// 	}
		// 	name := path[j+1:]

		// 	temp = ui.NewButton(name)
		// 	temp.OnClicked(func(*ui.Button) {
		// 		fmt.Println(path)
		// 		file.Dn(name)
		// 		label7.SetText("已下载")
		// 	})
		// 	box4.Append(temp, false)
		// 	label6.SetText("已上传")
		// })

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
