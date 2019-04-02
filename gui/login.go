package gui

import "log"

// Login creates the login window and collects infos.
func Login() {

	accessKey, secretKey, bucket, domain, zone, window, fileList :=
		LoginWindow()
	log.Println("Logined, close the login window.")

	FileWindow(
		accessKey, secretKey, bucket, domain, zone, window, fileList)

}
