package main

import (
	"ex09-10-11/internal/blackjack"
	"ex09-10-11/pkg/deck"
	"fmt"
	"slices"
)

type basicAI struct{}

func (ai *basicAI) Bet(_ bool) int {
	return 100
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

func (ai *basicAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	//TODO: Implement when countign cards
	// Noop (for now)
}

func main() {

	game := blackjack.New(blackjack.GameOptions{
		Decks:           4,
		Hands:           50000,
		BlackjackPayout: 1.5,
	})

	// winings := game.Play(blackjack.NewHumanAI())
	winings := game.Play(&basicAI{})
	fmt.Println(winings)
}
