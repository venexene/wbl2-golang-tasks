package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"github.com/venexene/minishell/minishell"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting directory: %s\n", err)
		}
		fmt.Printf("%s$ ", currentDir)

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		elems := strings.Split(input, " ")

		if len(elems) < 0 {
			continue
		}
		
		switch elems[0] {
		case "cd":
			if len(elems) < 2 {
				fmt.Println("cd: missing argument")
				continue
			}

			err = minishell.ChangeDirectory(elems[1])
			if err != nil {
				fmt.Printf("cd: %s\n", err)
			}

		case "pwd":
			err = minishell.PrintWorkingDirectory()
			if err != nil {
				fmt.Printf("pwd: %s\n", err)
			}

		case "echo":
			if len(elems) < 2 {
				fmt.Println("echo: missing argument")
				continue
			}
			
			minishell.Echo(elems[1:])
		
		case "kill":
			if len(elems) < 2 {
				fmt.Println("kill: missing argument")
				continue
			}
			
			minishell.Kill(elems[1])

		case "ps":
			err := minishell.ProcessStatus()
			if err != nil {
				fmt.Printf("ps: %s\n", err)
			}
		}
		
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}
}