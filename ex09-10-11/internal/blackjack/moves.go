package blackjack

import (
	"errors"
	"ex09-10-11/pkg/deck"
	"fmt"
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
	if g.state == DealerTurn {
		g.state++
		return nil
	}
	if g.state == PlayerTurn {
		g.handIdx++
		if g.handIdx >= len(g.player) {
			g.state++
		}
		return nil
	}
	return errors.New("invalid state")
}

// MoveDoubleDown doubles the player bet, hits and imediataly stands
func MoveDoubleDown(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("can only double with exactly two cards")
	}
	g.player[g.handIdx].bet *= 2
	MoveHit(g)
	return MoveStand(g)
}

func MoveSplit(g *Game) error {
	cHand := g.currentHand()
	if len(*cHand) != 2 {
		fmt.Println("can only split with exactly two cards in your hand")
		return errSplit
	}

	if (*cHand)[0].Value != (*cHand)[1].Value {
		fmt.Println("to split, both cards must have the same value")
		return errSplit
	}

	g.player = append(g.player, hand{
		cards: []deck.Card{(*cHand)[1]},
		bet:   g.player[g.handIdx].bet,
	})
	g.player[g.handIdx].cards = (*cHand)[:1]

	return nil
}
