package main

import (
	"fmt"
	"strings"
)

func main() {
	var input string
	fmt.Scanf("%s\n", &input)

	wordCount := 1
	for _, r := range input {
		// min, max := 'A', 'Z'
		// if r >= min && r <= max {
		// 	wordCount++
		// }
		str := string(r)

		if strings.ToUpper(str) == str {
			wordCount++
		}
	}
	fmt.Println(wordCount)
}
