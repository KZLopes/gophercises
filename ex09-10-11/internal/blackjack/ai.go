package blackjack

import (
	"ex09-10-11/pkg/deck"
	"fmt"
)

type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hands [][]deck.Card, dealer []deck.Card)
}

func NewHumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled.")
	}
	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
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

func (ai humanAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:")
	for _, hand := range hands {
		fmt.Println("\t", hand)
	}
	fmt.Println("Dealer:", dealer)
}
