package main

import (
	"ex09-10-11/internal/blackjack"
	"ex09-10-11/pkg/deck"
	"fmt"
	"slices"
)

type basicAI struct {
	score int
	seen  int
	decks int
}

func (ai *basicAI) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}

	trueScore := ai.score / ((ai.decks*51 - ai.seen) / 52)
	switch {
	case trueScore >= 14:
		return 10000
	case trueScore >= 8:
		return 500
	default:
		return 100
	}
}

func (ai *basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			shouldSplit := []deck.Value{deck.Ace, deck.Eight, deck.Nine}
			if slices.Contains(shouldSplit, hand[0].Value) {
				return blackjack.MoveSplit
			}
		}
		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDoubleDown
		}
	}

	dScore := blackjack.Score(dealer)
	if dScore >= 5 && dScore <= 6 {
		return blackjack.MoveStand
	}

	if score < 13 {
		return blackjack.MoveHit
	}

	return blackjack.MoveStand
}

func (ai *basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}
	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}

func (ai *basicAI) count(card deck.Card) {
	score := blackjack.Score(card)
	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	}
	ai.seen++
}

func main() {

	game := blackjack.New(blackjack.GameOptions{
		Decks:           4,
		Hands:           50000,
		BlackjackPayout: 1.5,
	})

	// winings := game.Play(blackjack.NewHumanAI())
	winings := game.Play(&basicAI{decks: 4})
	fmt.Println(winings)
}
