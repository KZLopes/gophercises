//go:generate stringer -type=Suit,Value

package deck

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"time"
)

type Suit uint8

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Joker
)

var suits = [...]Suit{Spades, Diamonds, Clubs, Hearts}

type Value uint8

const (
	_ Value = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minValue = Ace
	maxValue = King
)

type Card struct {
	Suit
	Value
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %s", c.Value.String(), c.Suit.String())
}

type deckOpt func([]Card) []Card

func New(opts ...deckOpt) []Card {
	var deck []Card
	for _, s := range suits {
		for v := minValue; v <= maxValue; v++ {
			deck = append(deck, Card{Suit: s, Value: v})
		}
	}

	for _, opt := range opts {
		deck = opt(deck)
	}
	return deck
}

func defaultSort(deck []Card) []Card {
	sort.Slice(deck, less(deck))
	return deck
}

func less(deck []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return rank(deck[i]) < rank(deck[j])
	}
}

func Sort(less func(deck []Card) func(i, j int) bool) deckOpt {
	return func(deck []Card) []Card {
		sort.Slice(deck, less(deck))
		return deck
	}
}

func rank(c Card) int {
	return int(c.Suit)*int(maxValue) + int(c.Value)
}

var seed = time.Now().Unix()

func Shuffle(deck []Card) []Card {
	// Getting a new source every time the methos is called
	r := rand.New(rand.NewSource(seed))

	for n := len(deck); n > 0; n-- {
		randIdx := r.Intn(n)
		deck[n-1], deck[randIdx] = deck[randIdx], deck[n-1]
	}
	return deck
	// Creating a New Slice (Easier to Understand Whats Going On)
	// ret := make([]Card, len(deck))
	// for i, j := range r.Perm(len(deck)) {
	// 	ret[i] = deck[j]
	// }
	// return ret
}

func Jokers(n int) deckOpt {
	return func(deck []Card) []Card {
		for i := 0; i < n; i++ {
			deck = append(deck, Card{
				Suit: Joker,
				// Value: Value(i),
			})
		}
		return deck
	}
}

// Filter could have functional options to filter specific Card.Value, Card.Suit or specific Card{}
func Filter(vals ...Value) deckOpt {
	return func(deck []Card) []Card {
		var ret = []Card{}
		for _, card := range deck {
			if !slices.Contains(vals, card.Value) {
				ret = append(ret, card)
			}
		}
		/* 		for _, v := range vals {
			for _, card := range deck {
				if card.Value != v {
					ret = append(ret, card)
				}
			}
		} */
		return ret
	}
}

func Decks(n int) deckOpt {
	return func(deck []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, deck...)
		}
		return ret
	}
}
