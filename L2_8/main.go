package main

import (
	"fmt"
	"log"
	"os"

	"github.com/venexene/myntp/myntp"
)

func main() {
	strTime, err := myntp.GetCurrentTime()
	if err != nil {
		log.Printf("Error getting time: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Current Time: %s\n", strTime)
}