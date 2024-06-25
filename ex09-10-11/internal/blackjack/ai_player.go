package blackjack

import (
	"ex09-10-11/pkg/deck"
	"fmt"
	"slices"
)

type PlayerAI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hands [][]deck.Card, dealer []deck.Card)
}

type humanPlayer struct{}

func NewHumanPlayer() humanPlayer {
	return humanPlayer{}
}

func (ai humanPlayer) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled.")
	}
	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanPlayer) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		str := "Pick your action: (h)it, (s)tand"
		if len(hand) == 2 {
			str += ", (d)oubledown"
			if hand[0].Value == hand[1].Value {
				str += ", s(p)lit"
			}
		}
		fmt.Println(str)

		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDoubleDown
		case "p":
			return MoveSplit
		default:
			fmt.Println("Invalid Action!")
		}
	}
}

func (ai humanPlayer) Results(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:")
	for _, hand := range hands {
		fmt.Println("\t", hand)
	}
	fmt.Println("Dealer:", dealer)
}

type defaultPlayer struct {
	score int
	seen  int
	decks int
}

func NewDefaultPlayer(decks int) defaultPlayer {
	return defaultPlayer{decks: decks}
}

func (ai *defaultPlayer) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}

	trueScore := ai.score / ((ai.decks*52 - ai.seen) / 52)
	switch {
	case trueScore >= 14:
		return 10000
	case trueScore >= 8:
		return 500
	default:
		return 100
	}
}

func (ai *defaultPlayer) Play(hand []deck.Card, dealer deck.Card) Move {
	score := Score(hand...)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			shouldSplit := []deck.Value{deck.Ace, deck.Eight, deck.Nine}
			if slices.Contains(shouldSplit, hand[0].Value) {
				return MoveSplit
			}
		}
		if (score == 10 || score == 11) && !Soft(hand...) {
			return MoveDoubleDown
		}
	}

	dScore := Score(dealer)
	if dScore >= 5 && dScore <= 6 {
		return MoveStand
	}

	if score < 13 {
		return MoveHit
	}

	return MoveStand
}

func (ai *defaultPlayer) Results(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}
	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}

func (ai *defaultPlayer) count(card deck.Card) {
	score := Score(card)
	switch {
	case score <= 6:
		ai.score++
	case score >= 10:
		ai.score--
	}
	ai.seen++
}
