// Package simplesort provides functionality for simple realization of sort from UNIX
package simplesort

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sort"
)



func getLinesFromFile(filepath string, unique bool) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s: %w", filepath, err)
	}
	defer file.Close()
	
	lines := make([][]string, 0)
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if unique {
			if seen[line] {
				continue
			}
			seen[line] = true
		}

		columns := strings.Fields(line)
		lines = append(lines, columns)
	}

	if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("Error reading file %s: %w", filepath, err)
    }

	return lines, nil
}



func writeLinesToFile(lines [][]string) error {
	output, err := os.Create("sortedFile.txt")
	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}
	defer output.Close()

	for i, line := range lines {
		output.WriteString(strings.Join(line, " "))
		if i != len(lines)-1 {
			output.WriteString("\n")
		}
	}

	return nil
}



func compareValues(a, b string, numeric bool) int {
	if numeric {
		n1, err1 := strconv.Atoi(a)
		n2, err2 := strconv.Atoi(b)
		if err1 == nil && err2 == nil {
			return n1 - n2
		}
	}

	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}



// CreateSortedFile sorts contents of file and writes it to another file
func CreateSortedFile(filepath string, column int, numeric bool, reverse bool, unique bool) error {
	lines, err := getLinesFromFile(filepath, unique)
	if err != nil {
		return fmt.Errorf("Failed to get lines from file %s: %w", filepath, err)
	}

	sort.Slice(lines, func(i, j int) bool {
		for col := column - 1; col < len(lines[i]) && col < len(lines[j]); col++ {
			compare := compareValues(lines[i][col], lines[j][col], numeric)
			if compare != 0 {
				if reverse {
					return compare > 0
				}
				return compare < 0
			}
		}

		if reverse {
			return i > j
		}
		return i < j
	})

	if err := writeLinesToFile(lines); err != nil {
		return fmt.Errorf("Error writing lines to file: %w", err)
	} 
	
	return nil
}