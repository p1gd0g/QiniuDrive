package gui

// Login creates the login window and collects infos.
func Login() {

	accessKey, secretKey, bucket, domain, zone, window, fileList :=
		LoginWindow()

	FileWindow(
		accessKey, secretKey, bucket, domain, zone, window, fileList)

}
