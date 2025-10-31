// Package telnet provides simple telnet
package telnet

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// RunTelnet runs telnet
func RunTelnet(host string, port string, timeout time.Duration) {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Fatalf("Connection timeout after %v", timeout)
		} else {
			log.Fatalf("Failed to connect: %s", err)
		}
	}
	defer conn.Close()

	log.Printf("Connected to %s", address)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(conn)
		for {
			data, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("Connection closed by server")
					return
				}
				log.Printf("Error reading from host: %s\n", err)
				return
			}
			os.Stdout.WriteString(data)
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text() + "\n"
			_, err := conn.Write([]byte(text))
			if err != nil {
				log.Printf("Write error: %s\n", err)
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from STDIN: %s\n", err)
		} else {
			log.Println("EOF received, closing connection")
		}
	}()

	wg.Wait()
	log.Println("Disconnected")
}
