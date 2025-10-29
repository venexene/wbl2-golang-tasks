package main

import (
	"flag"
	"github.com/venexene/cut/cut"
)

func main() {
	var (
		fields    = flag.String("f", "", "Номера полей для вывода")
		delimeter = flag.String("d", "\t", "Разделитель полей")
		separated = flag.Bool("s", false, "Выводить только строки с разделителем")
	)

	flag.Parse()

	cut.RunCut(*fields, *delimeter, *separated)
}
