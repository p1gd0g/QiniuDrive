package comm

import (
	"io"
	"log"
	"net/http"
	"os"
)

// Download the file using http.
func Download(
	name, domain string) (err error) {

	resp, err := http.Get("http://" +
		domain + "/" + name)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}
	log.Println("Download",
		name, "successfully.")

	return
}
