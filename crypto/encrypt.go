package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

// Encrypt the file with password and save it.
func Encrypt(pwd string, file string) {

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

	fmt.Println(out.Bytes())
	ioutil.WriteFile("testenc", out.Bytes(), 0644)

}

// FormatPwd formats the password so that the length is 32 Bytes.
func FormatPwd(pwd string) string {

	if len(pwd) > 32 {
		pwd = pwd[:32]
	}

	for len(pwd) < 32 {
		pwd += string(pwd[0])
	}

	return pwd
}
