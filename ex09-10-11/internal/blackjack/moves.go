package blackjack

import (
	"errors"
	"ex09-10-11/pkg/deck"
)

type Move func(*Game) error

func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

func MoveStand(g *Game) error {
	g.state++
	return nil
}

// MoveDoubleDown doubles the player bet, hits and imediataly stands
func MoveDoubleDown(g *Game) error {
	if len(g.player) != 2 {
		return errors.New("")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

func MoveSplit(g *Game) error { return nil }
