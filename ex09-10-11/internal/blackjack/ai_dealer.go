package blackjack

import "ex09-10-11/pkg/deck"

type dealerAI interface {
	pickMove(hand []deck.Card) Move
}

type defaultDealer struct{}

func (ai defaultDealer) pickMove(hand []deck.Card) Move {
	score := Score(hand...)
	if score < 16 || (score == 17 && Soft(hand...)) {
		return MoveHit
	}

	return MoveStand
}
