// Package anagrams provides functionality for finding anagrams in slice of strings
package anagrams

import (
	"fmt"
	"strings"
	"sort"
)

func createFrequencyKey(word string) string {
	runes := []rune(strings.ToLower(word))
	
	freqs := make(map[rune]int)
	for _, char := range runes {
		freqs[char]++
	}

	var chars []rune
	for char := range freqs {
		chars = append(chars, char)
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	key := ""
	for _, char := range chars {
		key += fmt.Sprintf("%c%d", char, freqs[char])
	}

	return key
}

func removeDuplicates(words []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}
	return result
}

// FindAnagrams finds anagrams and returns result as map
func FindAnagrams(words []string) map[string][]string {
	anagramGroups := make(map[string][]string)
	for _, word := range words {
		key := createFrequencyKey(word)
		anagramGroups[key] = append(anagramGroups[key], strings.ToLower(word))
	}

	result := make(map[string][]string)
	for _, words := range anagramGroups {
		if len(words) > 1 {
			unique := removeDuplicates(words)
			sort.Strings(unique)
			result[words[0]] = unique	
		}
	}


	return result
}