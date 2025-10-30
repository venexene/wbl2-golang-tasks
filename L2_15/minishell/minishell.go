package minishell

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func ChangeDirectory(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}
	return nil
}

func PrintWorkingDirectory() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(currentDir)
	return nil
}

func Echo(inputs []string) {
	fmt.Println(strings.Join(inputs, " "))
}

func Kill(pidStr string) error {
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

type Process struct {
	PID  int
	Name string
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

func ProcessStatus() error {
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