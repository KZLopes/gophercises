package main

import (
	"ex09-10-11/internal/blackjack"
	"fmt"
)

func main() {
	game := blackjack.New(blackjack.GameOptions{Hands: 2})
	winings := game.Play(blackjack.NewHumanAI())
	fmt.Println(winings)
}
