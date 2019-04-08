package main

import (
	"flag"
	"log"

	"github.com/p1gd0g/QiniuDrive/gui"
)

func main() {

	flag.String("ak", "", "Your access key.")
	flag.String("sk", "", "Your secret key.")
	flag.String("bk", "", "Your bucket name.")
	flag.String("dm", "", "Your domain name. Optional for Download.")
	flag.Int("zn", 0, "Your bucket zone: "+
		"1-Huadong, 2-Huabei, 3-Huanan, 4-Beimei. Optional for Upload.")

	flag.Parse()

	log.SetFlags(log.Lshortfile)

	gui.Start()
}
