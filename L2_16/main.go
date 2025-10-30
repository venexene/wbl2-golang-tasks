package main

import (
	"github.com/venexene/wget/wget"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 1 {
		log.Printf("Missing arguments")
		os.Exit(1)
	}

	link := os.Args[1]

	wget.MirrorPage(link)
}
