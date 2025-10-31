package main

import (
	"github.com/venexene/wget/wget"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing arguments")
	}

	link := os.Args[1]

	wget.MirrorPage(link)
}
