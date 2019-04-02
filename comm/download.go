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

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		return
	}

	resp, err := http.Get("http://" +
		domain + "/" + name)

	if err != nil {
		defer resp.Body.Close()
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}
	log.Println("Download",
		name, "successfully.")

	return
}
