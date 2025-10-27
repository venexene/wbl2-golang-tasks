// Package grep provides functionality for simple realization of grep
package grep

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func getLinesFromFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file %s: %w", filepath, err)
	}

	return lines, nil
}

// RunGrep sorts searches by pattern in file
func RunGrep(filepath string, pattern string, after int, before int, count bool, ignore bool, invert bool, fixed bool, number bool) error {
	lines, err := getLinesFromFile(filepath)
	if err != nil {
		return fmt.Errorf("Failed to get lines from file %s: %w", filepath, err)
	}

	var re *regexp.Regexp
	searchPattern := pattern
	if !fixed {
		if ignore {
			searchPattern = "(?i)" + pattern
		}
		re, err = regexp.Compile(searchPattern)
		if err != nil {
			return fmt.Errorf("Failed to compile pattern %s: %w", pattern, err)
		}
	}

	outputLines := make(map[int]bool)
	matchCount := 0
	for idx, line := range lines {
		lineToMatch := line
		if ignore && fixed {
			lineToMatch = strings.ToLower(lineToMatch)
		}

		var match bool
		if fixed {
			if ignore {
				match = strings.Contains(lineToMatch, strings.ToLower(pattern))
			} else {
				match = strings.Contains(lineToMatch, pattern)
			}
		} else {
			match = re.MatchString(lineToMatch)
		}

		if invert {
			match = !match
		}

		if match {
			matchCount++
			if !count {
				outputLines[idx] = true

				for i := 1; i <= after; i++ {
					if idx+i < len(lines) {
						outputLines[idx+i] = true
					}
				}

				for i := 1; i <= before; i++ {
					if idx-i >= 0 {
						outputLines[idx-i] = true
					}
				}
			}
		}
	}

	if count {
		fmt.Printf("%d\n", matchCount)
		return nil
	}

	for idx := 0; idx < len(lines); idx++ {
		if outputLines[idx] {
			if number {
				fmt.Printf("%d:%s\n", idx+1, lines[idx])
			} else {
				fmt.Println(lines[idx])
			}
		}
	}

	return nil
}
