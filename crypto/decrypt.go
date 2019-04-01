package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"io/ioutil"
	"log"

	"github.com/p1gd0g/QiniuDrive/tool"
)

// Decrypt the file with password and save it.
func Decrypt(pwd, file string) {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	log.Println("Read file successful.")

	bReader := bytes.NewReader(dat)

	block, err := aes.NewCipher([]byte(FormatPwd(pwd)))
	if err != nil {
		panic(err)
	}
	log.Println("Create new cipher successfully.")

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	var out bytes.Buffer

	writer := &cipher.StreamWriter{S: stream, W: &out}

	if _, err := io.Copy(writer, bReader); err != nil {
		panic(err)
	}
	log.Println("Copy the reader successfully.")

	ioutil.WriteFile("deced "+tool.GetFileName(file),
		out.Bytes(), 0644)
}
