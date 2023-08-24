package main

import (
	"flag"
)

var isServer *bool

func init() {
	isServer = flag.Bool("s", false, "help message for server")
	flag.Parse()
}

func main() {
	if *isServer {
		serve()
	} else {
		runClient()
	}
}
