package main

import (
	"fmt"

	"github.com/venexene/anagrams/anagrams"
)

func main () {
	resMap := anagrams.FindAnagrams([]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"})

	fmt.Println(resMap)
}