package comm

import (
	"log"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"
)

// Delete the file using qiniu's api.
func Delete(
	accessKey,
	secretKey,
	bucket,
	name string,
) (err error) {
	mac := auth.New(accessKey, secretKey)
	log.Println("mac created.")

	bucketManager := storage.NewBucketManager(
		mac, &storage.Config{})
	log.Println("bucketManager created.")

	return bucketManager.Delete(bucket, name)
}
