package main

import (
	"flag"
)

func main() {
	TYPE := flag.Bool("TYPE", false, "Type of server to start (mapper = true, reducer = false)")
	PORT := flag.String("PORT", "5001", "Port to start the server on")
	flag.Parse()

	StartWorker(*TYPE, *PORT)
}