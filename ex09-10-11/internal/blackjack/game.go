package blackjack

import (
	"ex09-10-11/pkg/deck"
	"fmt"
)

type state int8

const (
	Betting state = iota
	PlayerTurn
	DealerTurn
	Resolution
)

var (
	errBust error
)

type hand struct {
	cards []deck.Card
	bet   int
}

type GameOptions struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

type Game struct {
	nDecks          int
	nHands          int
	blackjackPayout float64

	state       state
	deck        []deck.Card
	reshuffleTh int

	player     []hand
	handIdx    int
	initialBet int
	balance    int

	dealer []deck.Card
}

func New(o GameOptions) Game {
	if o.Decks < 1 {
		o.Decks = 3
	}
	if o.Hands < 1 {
		o.Hands = 100
	}

	if o.BlackjackPayout < 1.0 {
		o.BlackjackPayout = 1.5
	}

	return Game{
		state:           Betting,
		nDecks:          o.Decks,
		nHands:          o.Hands,
		balance:         0,
		blackjackPayout: o.BlackjackPayout,
		reshuffleTh:     52 * o.Decks / 3,
	}
}

func (g *Game) Play(ai AI) int {
	g.deck = nil
	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) <= g.reshuffleTh {
			fmt.Println("The deck was just shuffled.")
			g.deck = deck.New(deck.Decks(g.nDecks), deck.Shuffle)
			shuffled = true
		}
		bet(g, ai, shuffled)
		deal(g)
		if Blackjack(g.dealer...) || Blackjack(g.player[g.handIdx].cards...) {
			endRound(g, ai)
			continue
		}

		for g.state == PlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			// hand := make([]hand, len(g.player))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				// Noop
			default:
				panic(err)
			}

		}

		dScore := Score(g.dealer...)
		for g.state == DealerTurn {
			if dScore <= 16 || (dScore == 17 && Soft(g.dealer...)) {
				MoveHit(g)
				dScore = Score(g.dealer...)
			} else {
				MoveStand(g)
			}
		}

		endRound(g, ai)
	}
	return g.balance
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case PlayerTurn:
		return &g.player[g.handIdx].cards
	case DealerTurn:
		return &g.dealer
	// case StHandOver:
	default:
		panic("Not a player`s turn")
	}
}

func bet(g *Game, ai AI, shuffled bool) {
	pBet := ai.Bet(shuffled)
	if pBet < 100 {
		fmt.Println("bet must be atleast 100")
		bet(g, ai, shuffled)
		return
	}
	g.initialBet = pBet
}

func deal(g *Game) {
	playerHand := make([]deck.Card, 0, 5)
	g.handIdx = 0
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card

	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		playerHand = append(playerHand, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.player = []hand{
		{
			cards: playerHand,
			bet:   g.initialBet,
		},
	}
	g.state = PlayerTurn
}

func endRound(g *Game, ai AI) {
	dScore, dBlackjack := Score(g.dealer...), Blackjack(g.dealer...)
	allPHands := make([][]deck.Card, len(g.player))
	for hIdx, hand := range g.player {
		cards := hand.cards
		allPHands[hIdx] = cards
		winnings := hand.bet
		pScore, pBlackjack := Score(hand.cards...), Blackjack(hand.cards...)
		switch {
		case pBlackjack && dBlackjack:
			winnings = 0
		case dBlackjack:
			winnings = -winnings
		case pBlackjack:
			winnings = int(float64(winnings) * g.blackjackPayout)
		case pScore > 21:
			winnings = -winnings
		case dScore > 21:
			// win
		case pScore > dScore:
			// win
		case dScore > pScore:
			winnings = -winnings
		case dScore == pScore:
			winnings = 0
		}
		g.balance += winnings
	}
	ai.Results(allPHands, g.dealer)
	g.player = nil
	g.dealer = nil
}

// Blackjack return true if a hand is a blackjack
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

// Score takes a hand of cards and returns the best blackjack score possible.
func Score(hand ...deck.Card) int {
	soft := minScore(hand...)
	if soft > 11 {
		return soft
	}
	for _, c := range hand {
		if c.Value == deck.Ace {
			return soft + 10
		}
	}
	return soft
}

// Soft retun true if the score of a hand is being counted with an Ace as 11 points.
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)

	return minScore != score
}

func minScore(hand ...deck.Card) int {
	score := 0

	for _, c := range hand {
		score += min(int(c.Value), 10)
	}

	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper to draw a single card
func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
