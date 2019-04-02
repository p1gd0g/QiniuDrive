package gui

import "log"

// Login creates the login window and inputs info.
func Login() {

	accessKey, secretKey, bucket, domain, zone, window, fileList :=
		NewLoginWindow()
	log.Println("Logined.")

	NewFileWindow(
		accessKey, secretKey, bucket, domain, zone, window, fileList)

}
