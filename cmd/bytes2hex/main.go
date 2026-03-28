package main

import (
	"log"

	bh "github.com/takanoriyanagitani/go-bytes2hex"
)

func main() {
	err := bh.BulkEncoderDefault.StdinToHexToStdout()
	if nil != err {
		log.Printf("%v\n", err)
	}
}
