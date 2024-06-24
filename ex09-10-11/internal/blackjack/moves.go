package blackjack

import "ex09-10-11/pkg/deck"

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
}

func MoveStand(g *Game) {
	g.state++
}

func MoveDoubleDown(g *Game) {}

func MoveSplit(g *Game) {}
