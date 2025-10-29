// Package cut provides functionality for simple realization of string cut by delimeter
package cut

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseFields(fields string) ([]int, error) {
	ranges := strings.Split(fields, ",")
	fieldVals := []int{}

	for _, r := range ranges {
		r = strings.TrimSpace(r)

		if num, err := strconv.Atoi(strings.TrimSpace(r)); err == nil {
			if num <= 0 {
				return nil, fmt.Errorf("Field number must be positive: %d", num)
			}
			fieldVals = append(fieldVals, num)
			continue
		}

		bounds := strings.Split(r, "-")

		if len(bounds) != 2 {
			return nil, fmt.Errorf("Invalid range format: %s", r)
		}

		lower, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
		if err != nil {
			return nil, fmt.Errorf("Invalid lower bound: %s", bounds[0])
		}

		upper, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
		if err != nil {
			return nil, fmt.Errorf("Invalid upper bound: %s", bounds[1])
		}

		if lower <= 0 || upper <= 0 {
			return nil, fmt.Errorf("Field numbers in range must be positive: %d-%d", lower, upper)
		}

		if lower > upper {
			return nil, fmt.Errorf("Invalid range: %d-%d", lower, upper)
		}

		for i := lower; i <= upper; i++ {
			fieldVals = append(fieldVals, i)
		}
	}

	return fieldVals, nil
}

// RunCut gets lines from STDIN and cuts them by delimeter
func RunCut(fields string, d string, s bool) {
	scanner := bufio.NewScanner(os.Stdin)

	fieldsNums, err := parseFields(fields)
	if err != nil {
		log.Fatalf("Error parsing fields: %v\n", err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, d)

		if s && len(cols) == 1 {
			continue
		}

		if len(cols) == 1 {
			fmt.Println(line)
			continue
		}

		outputFields := []string{}
		for _, fieldNum := range fieldsNums {
			idx := fieldNum - 1
			if idx >= 0 && idx < len(cols) {
				outputFields = append(outputFields, cols[idx])
			}
		}

		if len(outputFields) > 0 {
			fmt.Println((strings.Join(outputFields, d)))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}
}
