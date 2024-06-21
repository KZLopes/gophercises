package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Value: Ace, Suit: Spades})
	fmt.Println(Card{Value: Two, Suit: Hearts})
	fmt.Println(Card{Value: Seven, Suit: Clubs})
	fmt.Println(Card{Value: Jack, Suit: Diamonds})
	fmt.Println(Card{Value: Ace, Suit: Joker})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Spades
	// Two of Hearts
	// Seven of Clubs
	// Jack of Diamonds
	// Joker
	// Joker
}

// TestNew verifies that the New function creates a complete and correct deck of cards.
func TestNew(t *testing.T) {
	deck := New()

	// Check that the deck contains 52 cards
	expectedDeckSize := 52
	if len(deck) != expectedDeckSize {
		t.Errorf("Expected deck size of %d, but got %d", expectedDeckSize, len(deck))
	}

	// Verify that each suit has 13 cards
	expectedSuitsPerDeck := 4
	expectedCardsPerSuit := 13
	suitCount := make(map[Suit]int)
	valueCount := make(map[Suit]map[Value]int)
	for _, card := range deck {
		suitCount[card.Suit]++
		if valueCount[card.Suit] == nil {
			valueCount[card.Suit] = make(map[Value]int)
		}
		valueCount[card.Suit][card.Value]++
	}

	if len(suitCount) != expectedSuitsPerDeck {
		t.Errorf("Expected %d suits, but got %d", expectedSuitsPerDeck, len(suitCount))
	}

	for _, suit := range suits {
		if suitCount[suit] != expectedCardsPerSuit {
			t.Errorf("Expected 13 cards of suit %v, but got %d", suit, suitCount[suit])
		}
		for value := minValue; value <= maxValue; value++ {
			if valueCount[suit][value] != 1 {
				t.Errorf("Expected 1 card of value %v in suit %v, but got %d", value, suit, valueCount[suit][value])
			}
		}
	}
}

func TestDefaultSort(t *testing.T) {
	deck := New(defaultSort)

	if deck[0] != (Card{Suit: Spades, Value: Ace}) {
		t.Errorf("Expected first card to be Ace of Spades, but got %v", deck[0])
	}
}

func TestSort(t *testing.T) {
	deck := New(Sort(less))

	if deck[0] != (Card{Suit: Spades, Value: Ace}) {
		t.Errorf("Expected first card to be Ace of Spades, but got %v", deck[0])
	}
}

func TestShuffle(t *testing.T) {
	const shuffleTimes = 1000
	const minChangedPositions = 33

	sortedDeck := New()
	for i := 0; i < shuffleTimes; i++ {
		seed += int64(i)
		deck := New(Shuffle)

		var count uint8
		for idx, card := range deck {
			if card != sortedDeck[idx] {
				count++
			}
		}

		if count < 33 {
			t.Errorf("%d cards chenged position; expected atleast %d", count, minChangedPositions)
		}
	}

}

func TestJokers(t *testing.T) {
	testCases := []struct {
		input int
		want  int
		desc  string
	}{
		{
			input: 0,
			want:  52,
			desc:  "Zero Jokers",
		},
		{
			input: 1,
			want:  53,
			desc:  "One Jokers",
		},
		{
			input: 2,
			want:  54,
			desc:  "Two Jokers",
		},
		{
			input: 3,
			want:  55,
			desc:  "Three Jokers",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			deck := New(Jokers(tC.input))
			if len(deck) != tC.want {
				t.Errorf("got %d cards in the deck. Expected %d", len(deck), tC.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {

}

func TestDecks(t *testing.T) {
	testCases := []struct {
		decksAsked       int
		numCardsExpected int
		desc             string
	}{
		{
			decksAsked:       0,
			numCardsExpected: 0,
			desc:             "0 decks",
		},
		{
			decksAsked:       1,
			numCardsExpected: 52,
			desc:             "1 decks",
		},
		{
			decksAsked:       2,
			numCardsExpected: 52 * 2,
			desc:             "2 decks",
		},
		{
			decksAsked:       3,
			numCardsExpected: 52 * 3,
			desc:             "3 decks",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			deck := New(Decks(tC.decksAsked))
			if len(deck) != tC.numCardsExpected {
				t.Errorf("asked for %d decks got %d cards; Expected %d cards", tC.decksAsked, len(deck), tC.numCardsExpected)
			}
		})
	}
}
