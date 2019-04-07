package main

import (
	"flag"
	"log"

	"github.com/p1gd0g/QiniuDrive/gui"
)

func main() {

	flag.String("ak", "", "accessKey")
	flag.String("sk", "", "secretKey")
	flag.String("bk", "", "bucket")
	flag.String("dm", "", "domain")
	flag.Int("zn", 0, "zone")

	flag.Parse()

	log.SetFlags(log.Lshortfile)

	gui.Start()
}
