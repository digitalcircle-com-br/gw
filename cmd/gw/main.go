package main

import (
	"os"

	"github.com/digitalcircle-com-br/gw/lib/gw"
)

func main() {
	if len(os.Args) > 1 {
		os.Setenv("CONFIG", os.Args[1])
	}
	if os.Getenv("CONFIG") == "" {
		os.Setenv("CONFIG", "./gw.yaml")
	}
	gw.Init()
}
