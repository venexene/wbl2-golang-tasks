// Package minishell provides functionality for simple realization of shell
package minishell

import (
	"bufio"
	"fmt"
	"github.com/google/shlex"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func changeDirectory(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}
	return nil
}

func printWorkingDirectory() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(currentDir)
	return nil
}

func echo(inputs []string) {
	fmt.Println(strings.Join(inputs, " "))
}

func kill(pidStr string) error {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("Error converting pid string to int: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("Error finding process by pid: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("Error shutting process: %w", err)
	}

	return nil
}

func parseNameFromStatus(status string) string {
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Name:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return "Unknown"
}

func processStatus() error {
	files, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		commPath := filepath.Join("/proc", file.Name(), "comm")
		nameByte, err := os.ReadFile(commPath)

		var name string
		if err == nil {
			name = strings.TrimSpace(string(nameByte))
		} else {
			statusPath := filepath.Join("/proc", file.Name(), "status")
			statusByte, err := os.ReadFile(statusPath)
			if err == nil {
				name = parseNameFromStatus(string(statusByte))
			} else {
				name = "Unknown"
			}
		}

		fmt.Printf("Name: %s; PID: %d\n", name, pid)
	}
	return nil
}

func execute(inputStr string) error {
	parts, err := shlex.Split(inputStr)
	if err != nil {
		return err
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

// RunMinishell starts simple shell
func RunMinishell() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting directory: %s\n", err)
		}
		fmt.Printf("%s$ ", currentDir)

		if !scanner.Scan() {
			fmt.Println()
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

			err = changeDirectory(elems[1])
			if err != nil {
				fmt.Printf("cd: %s\n", err)
			}

		case "pwd":
			err = printWorkingDirectory()
			if err != nil {
				fmt.Printf("pwd: %s\n", err)
			}

		case "echo":
			if len(elems) < 2 {
				fmt.Println("echo: missing arguments")
				continue
			}

			echo(elems[1:])

		case "kill":
			if len(elems) < 2 {
				fmt.Println("kill: missing arguments")
				continue
			}

			kill(elems[1])
			if err != nil {
				fmt.Printf("kill: %s\n", err)
			}

		case "ps":
			err := processStatus()
			if err != nil {
				fmt.Printf("ps: %s\n", err)
			}

		case "exec":
			if len(elems) < 2 {
				fmt.Println("exec: missing arguments")
				continue
			}

			err := execute(strings.Join(elems[1:], " "))
			if err != nil {
				fmt.Printf("exec: %s\n", err)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}
}
