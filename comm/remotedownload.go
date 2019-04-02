package comm

import (
	"log"
	"strings"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"
)

// RemoteDownload fetches the url using qiniu's api.
func RemoteDownload(
	accessKey,
	secretKey,
	bucket,
	url string,
) (err error) {

	fileName :=
		url[strings.LastIndex(url, "/")+1:]

	mac := auth.New(accessKey, secretKey)
	log.Println("mac created.")

	bucketManager := storage.NewBucketManager(mac, &storage.Config{})
	log.Println("bucketManager created.")

	_, err = bucketManager.Fetch(url,
		bucket, fileName)
	return
}
