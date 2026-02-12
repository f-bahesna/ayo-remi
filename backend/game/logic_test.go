package game

import (
	"fmt"
	"testing"
	"remi-game/models"
)

func TestIsValidSet(t *testing.T) {
	tests := []struct {
		name     string
		cards    []models.Card
		expected bool
	}{
     {
			name: "Valid Run A-2-3 (Low Ace)",
			cards: []models.Card{
				{Suit: models.Hearts, Rank: models.Ace},
				{Suit: models.Hearts, Rank: 2},
				{Suit: models.Hearts, Rank: 3},
			},
			expected: true,
		},
        {
			name: "Valid Run Q-K-A (High Ace)",
			cards: []models.Card{
				{Suit: models.Clubs, Rank: models.Queen},
				{Suit: models.Clubs, Rank: models.King},
				{Suit: models.Clubs, Rank: models.Ace},
			},
			expected: true,
		},
        {
			name: "Invalid Run K-A-2 (Wrapping)",
			cards: []models.Card{
				{Suit: models.Diamonds, Rank: models.King},
				{Suit: models.Diamonds, Rank: models.Ace},
				{Suit: models.Diamonds, Rank: 2},
			},
			expected: false,
		},
        {
            name: "Valid Generic Set 7-7-7 (Different Suits)",
            cards: []models.Card{
                {Suit: models.Hearts, Rank: 7},
                {Suit: models.Spades, Rank: 7},
                {Suit: models.Diamonds, Rank: 7},
            },
            expected: true,
        },
        {
            name: "Invalid Ace Set (Duplicate Suits)",
            cards: []models.Card{
                {Suit: models.Hearts, Rank: models.Ace},
                {Suit: models.Hearts, Rank: models.Ace}, // Duplicate Suit
                {Suit: models.Diamonds, Rank: models.Ace},
            },
            expected: false, // Strict Unique Suit Rule
        },
		{
			name: "Valid Run 1-2-3 Same Suit (Explicit 1)",
			cards: []models.Card{
				{Suit: models.Hearts, Rank: 1},
				{Suit: models.Hearts, Rank: 2},
				{Suit: models.Hearts, Rank: 3},
			},
			expected: true,
		},
		{
			name: "Invalid Run 1-2-3 Different Suit",
			cards: []models.Card{
				{Suit: models.Hearts, Rank: 1},
				{Suit: models.Spades, Rank: 2},
				{Suit: models.Hearts, Rank: 3},
			},
			expected: false,
		},
		{
			name: "Valid Run J-Q-K Same Suit",
			cards: []models.Card{
				{Suit: models.Clubs, Rank: models.Jack},
				{Suit: models.Clubs, Rank: models.Queen},
				{Suit: models.Clubs, Rank: models.King},
			},
			expected: true,
		},
		{
			name: "Valid Triple A (Any Suit)",
			cards: []models.Card{
				{Suit: models.Hearts, Rank: models.Ace},
				{Suit: models.Spades, Rank: models.Ace},
				{Suit: models.Diamonds, Rank: models.Ace},
			},
			expected: true,
		},
		{
			name: "Valid Run with Joker",
			cards: []models.Card{
				{Suit: models.Hearts, Rank: 1},
				{Suit: models.Joker, Rank: models.JokerRank},
				{Suit: models.Hearts, Rank: 3},
			},
			expected: true,
		},
        {
			name: "Invalid Set (Random cards)",
			cards: []models.Card{
				{Suit: models.Hearts, Rank: 5},
				{Suit: models.Spades, Rank: 9},
				{Suit: models.Diamonds, Rank: 2},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSet(tt.cards); got != tt.expected {
				t.Errorf("IsValidSet() = %v, want %v for %v", got, tt.expected, tt.name)
			}
		})
	}
}

func TestNewDeck_Uniqueness(t *testing.T) {
    deck := NewDeck()
    
    // 1. Check Count (52 + 2 Jokers = 54)
    if len(deck) != 54 {
        t.Errorf("Expected 54 cards, got %d", len(deck))
    }
    
    // 2. Check Uniqueness (Suit + Rank)
    seen := make(map[string]bool)
    counts := make(map[models.Suit]int)

    for _, c := range deck {
        if c.Suit == models.Joker {
            continue 
        }
        
        // ID is unique by definition (UUID), so check Suit+Rank
        key := string(c.Suit) + "-" + fmt.Sprintf("%d", c.Rank)
        
        if seen[key] {
            t.Errorf("Duplicate card found: %s", key)
        }
        seen[key] = true
        counts[c.Suit]++
    }
    
    // 3. Confirm Suit Counts (13 of each)
    for _, s := range []models.Suit{models.Spades, models.Hearts, models.Diamonds, models.Clubs} {
        if counts[s] != 13 {
             t.Errorf("Expected 13 cards for suit %s, got %d", s, counts[s])
        }
    }
}

func TestIsWinningHand(t *testing.T) {
    tests := []struct {
        name     string
        cards    []models.Card
        expected bool
    }{
        {
            name: "Simple Win (Run + Discard)",
            cards: []models.Card{
                {Suit: models.Hearts, Rank: 1}, // A
                {Suit: models.Hearts, Rank: 2},
                {Suit: models.Hearts, Rank: 3},
                {Suit: models.Spades, Rank: models.King}, // Discard
            },
            expected: true,
        },
        {
            name: "Win with Joker (Run + Joker + Discard)",
            cards: []models.Card{
                {Suit: models.Hearts, Rank: 8},
                {Suit: models.Hearts, Rank: 9},
                {Suit: models.Joker, Rank: models.JokerRank}, // Acts as 10 or 7
                {Suit: models.Clubs, Rank: models.King}, // Discard (User scenario)
            },
            expected: true,
        },
        {
            name: "Double Set Win (3+3 + Discard)",
            cards: []models.Card{
                {Suit: models.Hearts, Rank: 2}, {Suit: models.Hearts, Rank: 3}, {Suit: models.Hearts, Rank: 4},
                {Suit: models.Spades, Rank: 5}, {Suit: models.Spades, Rank: 6}, {Suit: models.Spades, Rank: 7},
                {Suit: models.Clubs, Rank: 10}, // Discard
            },
            expected: true,
        },
        {
            name: "Losing Hand (Garbage)",
            cards: []models.Card{
                {Suit: models.Hearts, Rank: 2},
                {Suit: models.Spades, Rank: 5},
                {Suit: models.Diamonds, Rank: 9},
                {Suit: models.Clubs, Rank: models.King},
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := IsWinningHand(tt.cards); got != tt.expected {
                t.Errorf("IsWinningHand() = %v, want %v", got, tt.expected)
            }
        })
    }
}
