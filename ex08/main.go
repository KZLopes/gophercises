package main

import (
	"bytes"
	"fmt"
	"regexp"
)

var pNumbers = []string{"1234567890", "123 456 7891", "(123) 456 7892", "(123) 456-7893", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892"}

func main() {
	var normalizedNumbers []string

	for _, n := range pNumbers {
		normalized := normalize(n)
		normalizedNumbers = append(normalizedNumbers, string(normalized))
	}

	fmt.Println(normalizedNumbers)

}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func normalizeWithRegex(phone string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phone, "")
}
