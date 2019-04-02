package comm

import (
	"context"
	"log"

	"github.com/p1gd0g/QiniuDrive/tool"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/storage"
)

// Upload the file using qiniu's api.
func Upload(
	accessKey,
	secretKey,
	bucket,
	file string,
	zone int) (err error) {

	fileName := tool.GetFileName(file)

	mac := auth.New(accessKey, secretKey)
	log.Println("mac created.")

	cfg := storage.Config{}

	switch zone {
	case 0:
		cfg.Zone = &storage.ZoneHuadong
		log.Println("Zone: Huadong.")
	case 1:
		cfg.Zone = &storage.ZoneHuabei
		log.Println("Zone: Huabei.")
	case 2:
		cfg.Zone = &storage.ZoneHuanan
		log.Println("Zone: Huanan.")
	case 3:
		cfg.Zone = &storage.ZoneBeimei
		log.Println("Zone: Beimei.")
	default:
		log.Println("No zone!")
	}

	putPolicy := storage.PutPolicy{
		Scope: bucket + ":" + fileName,
	}
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err = formUploader.PutFile(context.Background(),
		&ret, upToken, fileName, file, &storage.PutExtra{})

	return
}
