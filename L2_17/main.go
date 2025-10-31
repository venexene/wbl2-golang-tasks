package main

import (
	"flag"
	"log"
	"time"

	"github.com/venexene/telnet/telnet"
)

func main() {
	timeoutFlag := flag.Uint("timeout", 10, "указать таймаут подключения")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalf("Missing arguments")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	timeout := time.Duration(*timeoutFlag) * time.Second

	telnet.RunTelnet(host, port, timeout)
}
