package main

import (
	"flag"
	"fmt"
	"github.com/venexene/grep/grep"
	"log"
	"os"
	"strings"
)

func parseArguments(args []string) []string {
	var parsed []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && len(arg) > 2 {
			if !strings.ContainsAny(arg, "0123456789") {
				for _, char := range arg[1:] {
					parsed = append(parsed, "-"+string(char))
				}
			} else {
				parsed = append(parsed, arg)
			}
		} else {
			parsed = append(parsed, arg)
		}
	}
	return parsed
}

func main() {
	var (
		after   = flag.Int("A", 0, "Выводить N строк после")
		before  = flag.Int("B", 0, "Выводить N строк до")
		context = flag.Int("C", 0, "Выводить N строк контекста вокруг")
		count   = flag.Bool("c", false, "Выводить только количество")
		ignore  = flag.Bool("i", false, "Игнорировать регистр")
		invert  = flag.Bool("v", false, "Инвертировать фильтр")
		fixed   = flag.Bool("F", false, "Воспринимать шаблон как фиксированную строку")
		number  = flag.Bool("n", false, "Выводить номер строки")
	)

	parsedArgs := parseArguments(os.Args[1:])

	oldArgs := os.Args
	os.Args = append([]string{oldArgs[0]}, parsedArgs...)
	flag.Parse()
	os.Args = oldArgs

	if len(flag.Args()) < 1 {
		fmt.Println("Введите название файла и шаблон!")
		os.Exit(1)
	}
	pattern := flag.Arg(0)
	filepath := flag.Arg(1)

	if *context > 0 {
		*after = *context
		*before = *context
	}

	if err := grep.RunGrep(filepath, pattern, *after, *before, *count, *ignore, *invert, *fixed, *number); err != nil {
		log.Printf("Error creatung sorted file: %v", err)
	}
}
