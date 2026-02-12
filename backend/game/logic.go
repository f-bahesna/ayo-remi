package game

import (
	"math/rand"
	"sort"
	"time"

	"remi-game/models"
    "github.com/google/uuid"
)

func NewDeck() []models.Card {
	suits := []models.Suit{models.Spades, models.Hearts, models.Diamonds, models.Clubs}
	var deck []models.Card

	// Add 2-10, J, Q, K, A for each suit
	for _, suit := range suits {
		// 2-10 (Fix: was 1-9, causing missing 10s and duplicate '1' Aces)
		for r := 2; r <= 10; r++ {
			deck = append(deck, models.Card{Suit: suit, Rank: models.Rank(r), ID: uuid.New().String()})
		}
		// J, Q, K, A
		ranks := []models.Rank{models.Jack, models.Queen, models.King, models.Ace}
		for _, r := range ranks {
			deck = append(deck, models.Card{Suit: suit, Rank: r, ID: uuid.New().String()})
		}
	}

	// Add 2 Jokers
	deck = append(deck, models.Card{Suit: models.Joker, Rank: models.JokerRank, ID: uuid.New().String()})
	deck = append(deck, models.Card{Suit: models.Joker, Rank: models.JokerRank, ID: uuid.New().String()})

	return deck
}

func Shuffle(deck []models.Card) []models.Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]models.Card, len(deck))
	copy(shuffled, deck)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

// IsValidSet checks if a group of cards forms a valid run or set.
// Rules:
// 1. Run (Sequence): Same suit, consecutive ranks (min 3). 
//    - Supports High Ace (Q-K-A) and Low Ace (A-2-3).
//    - K-A-2 is NOT valid (no wrapping).
// 2. Set (Group): Same rank, different suits (min 3).
// Jokers are treated as wildcards.
func IsValidSet(cards []models.Card) bool {
	if len(cards) < 3 {
		return false
	}

	// Separate Jokers from regular cards
	var jokers []models.Card
	var regulars []models.Card
	for _, c := range cards {
		if c.Suit == models.Joker {
			jokers = append(jokers, c)
		} else {
			regulars = append(regulars, c)
		}
	}

    numJokers := len(jokers)
    
    // If all cards are jokers, valid if >= 3
    if len(regulars) == 0 {
        return true 
    }

	// --- CHECK 1: SET (Group of Same Rank) ---
    // All regular cards must have the same rank.
    firstRank := regulars[0].Rank
    isSet := true
    seenSuits := make(map[models.Suit]bool)
    
    for _, c := range regulars {
        if c.Rank != firstRank {
            isSet = false
            break
        }
        // Strict Rule: Sets must have unique suits (e.g., A♠ A♥ A♣).
        // Standard Rummy does not allow duplicate suits in a set (single deck).
        if seenSuits[c.Suit] {
            isSet = false
            break
        }
        seenSuits[c.Suit] = true
    }
    
    if isSet {
        // Valid Set found (Group of same rank, unique suits)
        return true
    }

    // --- CHECK 2: RUN (Sequence of Same Suit) ---
    // All regular cards must have the same suit.
    firstSuit := regulars[0].Suit
    isRunSuit := true
    for _, c := range regulars {
        if c.Suit != firstSuit {
            isRunSuit = false
            break
        }
    }
    
    if !isRunSuit {
        return false // Neither a Set nor a Run
    }

    // It is a sequence of same suit. Check for consecutive ranks.
    // We need to handle Jokers filling gaps.
    
    // Check High Ace (A=14)
	sort.Slice(regulars, func(i, j int) bool {
		return regulars[i].Rank < regulars[j].Rank
	})
    
    if checkConsecutive(regulars, numJokers) {
        return true
    }

    // Check Low Ace (A=1)
    // If we have Aces, try treating them as Rank 1
    hasAce := false
    var lowAceRegulars []models.Card
    for _, c := range regulars {
        if c.Rank == models.Ace {
            hasAce = true
            lowAceRegulars = append(lowAceRegulars, models.Card{
                Suit: c.Suit,
                Rank: 1, // Treat as 1
                ID: c.ID,
            })
        } else {
             lowAceRegulars = append(lowAceRegulars, c)
        }
    }

    if hasAce {
        sort.Slice(lowAceRegulars, func(i, j int) bool {
            return lowAceRegulars[i].Rank < lowAceRegulars[j].Rank
        })
        if checkConsecutive(lowAceRegulars, numJokers) {
            return true
        }
    }

	return false
}

// Helper to check consecutive ranks with jokers
func checkConsecutive(sortedCards []models.Card, numJokers int) bool {
    gaps := 0
    for i := 0; i < len(sortedCards)-1; i++ {
        diff := int(sortedCards[i+1].Rank - sortedCards[i].Rank)
        if diff <= 0 {
            // Duplicate rank in a run is invalid (e.g. 5-5-6)
            return false
        }
        gaps += diff - 1
    }
    return gaps <= numJokers
}

// CanFormSetWith checks if adding newCard to hand can form at least one valid set
func CanFormSetWith(hand []models.Card, newCard models.Card) bool {
    return CanFormSetWithMultiple(hand, []models.Card{newCard})
}

// CanFormSetWithMultiple checks if ALL newCards can be part of a SINGLE valid set
// when combined with cards from the hand.
func CanFormSetWithMultiple(hand []models.Card, newCards []models.Card) bool {
    // Combined set must be valid and include ALL newCards.
    // If newCards alone are valid, return true.
    if IsValidSet(newCards) {
        return true
    }

    // Try combining newCards with every subset of hand.
    // Since hand size is small (max ~13), 2^13 is ~8192 iterations, which is fine.
    n := len(hand)
    count := 1 << n
    
    for i := 0; i < count; i++ {
        var subset []models.Card
        for j := 0; j < n; j++ {
            if (i & (1 << j)) > 0 {
                subset = append(subset, hand[j])
            }
        }
        
        // Combine subset with newCards
        // We create a new slice to avoid modifying underlying arrays if any
        testSet := make([]models.Card, 0, len(subset)+len(newCards))
        testSet = append(testSet, subset...)
        testSet = append(testSet, newCards...)
        
        if IsValidSet(testSet) {
            return true
        }
    }
    
    return false
}

// IsWinningHand checks if the hand can be fully formed into valid sets
// with at most 1 card remaining (the discard).
func IsWinningHand(hand []models.Card) bool {
    // If hand is empty or 1 card, it's a win (trivial)
    if len(hand) <= 1 {
        return true
    }
    
    // We need to check if there exists a card 'C' in Hand such that
    // Hand - {C} can be partitioned into valid sets.
    // Optimization: If len(hand) % 3 == 0, we might look for 0 discard?
    // In Rummy, you typically discard to win. So len(hand) should be 3k+1.
    // But let's be flexible: Allow 0 or 1 leftovers.
    
    // Try finding a perfect partition (0 leftovers)
    if canPartition(hand) {
        return true
    }
    
    // Try removing 1 card and partitioning the rest
    for i := range hand {
        // Create subset without card i
        // copy to avoid mutation issues
        subset := make([]models.Card, 0, len(hand)-1)
        subset = append(subset, hand[:i]...)
        subset = append(subset, hand[i+1:]...)
        
        if canPartition(subset) {
            return true
        }
    }
    
    return false
}

// canPartition recursively checks if cards can be split into valid sets
func canPartition(cards []models.Card) bool {
    if len(cards) == 0 {
        return true
    }
    if len(cards) < 3 {
        return false
    }
    
    // Try to find a valid set including the first card
    // We must use the first card to avoid trying same sets multiple times
    first := cards[0]
    rest := cards[1:]
    
    // We need to choose 2 or more cards from 'rest' to form a set with 'first'
    // This is the subset sum problem variant.
    // N is small (max ~13).
    
    n := len(rest)
    // Iterate all subsets of rest of size >= 2
    // To do this efficiently:
    // We can just iterate all subsets, check if isValid(subset + first), then recurse on remainder.
    
    maxSubset := 1 << n
    for i := 0; i < maxSubset; i++ {
        var potentialSet []models.Card
         var remaining []models.Card
         
        potentialSet = append(potentialSet, first)
        
        // Build potential set and remaining list
        for j := 0; j < n; j++ {
            if (i & (1 << j)) > 0 {
                potentialSet = append(potentialSet, rest[j])
            } else {
                remaining = append(remaining, rest[j])
            }
        }
        
        if len(potentialSet) >= 3 && IsValidSet(potentialSet) {
            if canPartition(remaining) {
                return true
            }
        }
    }
    
    return false
}
