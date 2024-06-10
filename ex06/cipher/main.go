package main

import (
	"fmt"
	"strings"
)

var alphabetLower = "abcdefghijklmnopqrstuvwxyz"
var alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func rotate(s rune, delta int, key []rune) rune {
	idx := strings.IndexRune(string(key), s)
	if idx < 0 {
		panic("idx < 0")
	}
	idx = (idx + delta) % len(key)
	return key[idx]
}

func main() {
	var length, secret int
	var message string
	fmt.Scanf("%d\n", &length)
	fmt.Scanf("%s\n", &message)
	fmt.Scanf("%d\n", &secret)

	var ret string

	for _, r := range message {
		switch {
		case strings.IndexRune(alphabetLower, r) >= 0:
			ret += string(rotate(r, secret, []rune(alphabetLower)))
		case strings.IndexRune(alphabetUpper, r) >= 0:
			ret += string(rotate(r, secret, []rune(alphabetUpper)))
		default:
			ret += string(r)
		}
	}

	fmt.Println(ret)
}
