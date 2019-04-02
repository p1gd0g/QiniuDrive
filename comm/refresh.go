package comm

import (
	"log"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"
)

// Refresh refreshes the file list.
func Refresh(
	accessKey,
	secretKey,
	bucket string,
) (
	list []storage.ListItem,
	err error,
) {
	mac := auth.New(accessKey, secretKey)
	log.Println("mac created.")

	bucketManager := storage.NewBucketManager(mac, &storage.Config{})
	log.Println("bucketManager created.")

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

		list = append(list, entries...)

		if hashNext {
			marker = nextMarker
		} else {
			break
		}

	}
	log.Println("Refreshed the file list.")
	return
}
