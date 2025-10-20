package main

import (
	"fmt"
	"log"

	"github.com/venexene/unpack/unpack"
)

func main() {
	newStr, err := unpack.Unpack(`av\55b4c10`)
	if err != nil {
		log.Printf("Error unpacking string: %v", err)
	}
	fmt.Println(newStr)
}