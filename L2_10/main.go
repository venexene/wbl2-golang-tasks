package main

import (
	"flag"
	"os"
	"fmt"
	"strings"
	"log"
	"github.com/venexene/sort/simplesort"
)

func parseArguments(args []string) []string {
	var parsed []string
	
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			for _, char := range arg[1:] {
				parsed = append(parsed, "-"+string(char))
			}
		} else {
			parsed = append(parsed, arg)
		}
	}

	return parsed
}

func main() {
	var (
		column  	  = flag.Int("k", 0, "Сортировать по столбцу")
		numVal  	  = flag.Bool("n", false, "Сортировать по числовому значению")
		reverse 	  = flag.Bool("r", false, "Сортировать в обратном порядке")
		unique  	  = flag.Bool("u", false, "Выводить только уникальные")
	)

	parsedArgs := parseArguments(os.Args[1:])

	oldArgs := os.Args
	os.Args = append([]string{oldArgs[0]}, parsedArgs...)
	flag.Parse()
	os.Args = oldArgs

	if len(flag.Args()) < 1 {
		fmt.Println("Введите название файла!")
		os.Exit(1)
	}
	filepath := flag.Arg(0)

	if err := simplesort.CreateSortedFile(filepath, *column, *numVal, *reverse, *unique); err != nil {
		log.Printf("Error creatung sorted file: %v", err)
	}
}