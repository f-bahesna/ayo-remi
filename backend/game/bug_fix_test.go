package game

import (
    "testing"
    "time"
    "remi-game/models"
)

// Bug 1: Joker must work as wildcard in any position
func TestBug1_JokerAsWildcard(t *testing.T) {
    store := &MockStore{}
    
    // Scenario: Player has [2S, Joker, 4S, 10H, JH, QH, KH]
    // They want to play [2S, Joker, 4S] as a valid run (Joker = 3S)
    gm := NewGame(store)
    p1 := models.Player{ID: "p1", Name: "P1", Hand: []models.Card{
        {Suit: models.Spades, Rank: 2, ID: "h1"},
        {Suit: models.Joker, Rank: models.JokerRank, ID: "joker1"},
        {Suit: models.Spades, Rank: 4, ID: "h3"},
        {Suit: models.Hearts, Rank: 10, ID: "h4"},
        {Suit: models.Hearts, Rank: 11, ID: "h5"},
        {Suit: models.Hearts, Rank: 12, ID: "h6"},
        {Suit: models.Hearts, Rank: 13, ID: "h7"},
    }}
    gm.Game = models.Game{
        ID:     "test",
        Status: models.StateInProgress,
        Players: []models.Player{p1},
        CurrentTurnPlayer: 0,
        TurnPhase: models.PhasePlay,
        TableSets: make([][]models.Card, 0),
    }

    // Play the set [2S, Joker, 4S]
    err := gm.PlaySet("p1", []models.Card{
        {Suit: models.Spades, Rank: 2, ID: "h1"},
        {Suit: models.Joker, Rank: models.JokerRank, ID: "joker1"},
        {Suit: models.Spades, Rank: 4, ID: "h3"},
    })
    if err != nil {
        t.Fatalf("Bug 1 STILL BROKEN: PlaySet [2S, Joker, 4S] failed: %v", err)
    }

    // Hand should now have 4 cards
    if len(gm.Game.Players[0].Hand) != 4 {
        t.Errorf("Hand should have 4 cards, got %d", len(gm.Game.Players[0].Hand))
    }

    // Table should have 1 set
    if len(gm.Game.Players[0].PlayedSets) != 1 {
        t.Errorf("Should have 1 played set, got %d", len(gm.Game.Players[0].PlayedSets))
    }

    t.Log("Bug 1 FIXED: Joker works as wildcard in [2S, Joker, 4S]")
}

// Bug 1b: Joker via DrawFromPile
func TestBug1b_JokerFromPile(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)

    // Hand has [2S, 4S, other cards...], Joker is in pile (last 3)
    p1 := models.Player{ID: "p1", Name: "P1", Hand: []models.Card{
        {Suit: models.Spades, Rank: 2, ID: "h1"},
        {Suit: models.Spades, Rank: 4, ID: "h2"},
        {Suit: models.Hearts, Rank: 10, ID: "h3"},
        {Suit: models.Hearts, Rank: 11, ID: "h4"},
        {Suit: models.Hearts, Rank: 12, ID: "h5"},
    }}
    gm.Game = models.Game{
        ID:     "test",
        Status: models.StateInProgress,
        Players: []models.Player{p1},
        CurrentTurnPlayer: 0,
        TurnPhase: models.PhaseDraw,
        Pile: []models.Card{
            {Suit: models.Joker, Rank: models.JokerRank, ID: "joker_pile"},
        },
        TableSets: make([][]models.Card, 0),
    }

    // DrawFromPile: handCards [2S, 4S] + pileCard Joker = Run [2,3,4]S
    err := gm.DrawFromPile("p1", []string{"h1", "h2"}, "joker_pile")
    if err != nil {
        t.Fatalf("Bug 1b STILL BROKEN: DrawFromPile with Joker failed: %v", err)
    }

    if len(gm.Game.Players[0].Hand) != 3 {
        t.Errorf("Hand should have 3 cards, got %d", len(gm.Game.Players[0].Hand))
    }
    if gm.Game.TurnPhase != models.PhasePlay {
        t.Errorf("Phase should be PLAY after DrawFromPile, got %s", gm.Game.TurnPhase)
    }

    t.Log("Bug 1b FIXED: Joker from pile works with DrawFromPile")
}

// Bug 2: Must be able to discard after DrawFromPile
func TestBug2_DiscardAfterDrawFromPile(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)

    p1 := models.Player{ID: "p1", Name: "P1", Hand: []models.Card{
        {Suit: models.Spades, Rank: 4, ID: "h1"},
        {Suit: models.Spades, Rank: 6, ID: "h2"},
        {Suit: models.Hearts, Rank: 10, ID: "h3"},
        {Suit: models.Hearts, Rank: 11, ID: "h4"},
        {Suit: models.Hearts, Rank: 12, ID: "h5"},
    }}
    p2 := models.Player{ID: "p2", Name: "P2", Hand: []models.Card{
        {Suit: models.Clubs, Rank: 2, ID: "p2h1"},
    }}

    gm.Game = models.Game{
        ID:     "test",
        Status: models.StateInProgress,
        Players: []models.Player{p1, p2, 
            {ID: "p3", Name: "Bot1", Hand: []models.Card{}},
            {ID: "p4", Name: "Bot2", Hand: []models.Card{}},
        },
        CurrentTurnPlayer: 0,
        TurnPhase: models.PhaseDraw,
        Pile: []models.Card{
            {Suit: models.Spades, Rank: 5, ID: "pile1"},
        },
        TableSets: make([][]models.Card, 0),
    }

    // Step 1: DrawFromPile (4S + 6S + 5S from pile = run)
    err := gm.DrawFromPile("p1", []string{"h1", "h2"}, "pile1")
    if err != nil {
        t.Fatalf("DrawFromPile failed: %v", err)
    }

    // After DrawFromPile: phase should be PLAY
    if gm.Game.TurnPhase != models.PhasePlay {
        t.Fatalf("Phase should be PLAY after DrawFromPile, got %s", gm.Game.TurnPhase)
    }

    // Step 2: Try to discard - THIS IS THE BUG TEST
    // Hand still has [10H, JH, QH] (3 cards left)
    err = gm.DiscardCard("p1", "h3") // Discard 10H
    if err != nil {
        t.Fatalf("Bug 2 STILL BROKEN: DiscardCard after DrawFromPile failed: %v", err)
    }

    // Verify turn moved to next player
    if gm.Game.CurrentTurnPlayer != 1 {
        t.Errorf("Turn should move to player 1, got %d", gm.Game.CurrentTurnPlayer)
    }
    if gm.Game.TurnPhase != models.PhaseDraw {
        t.Errorf("Next player phase should be DRAW, got %s", gm.Game.TurnPhase)
    }

    // Hand should now have 2 cards
    if len(gm.Game.Players[0].Hand) != 2 {
        t.Errorf("P1 hand should have 2 cards, got %d", len(gm.Game.Players[0].Hand))
    }

    t.Log("Bug 2 FIXED: Discard works correctly after DrawFromPile")
}

// Bug 3: RestartGame/StartGame must not deadlock when called from a context
// that already holds gm.Mutex, which is exactly how ws/client.go's
// MapMessageToGameAction dispatches every message (it locks gm.Mutex once for
// the whole switch statement, then calls into the game package). The locking
// RestartGame/StartGame methods re-lock the same non-reentrant sync.Mutex,
// which previously deadlocked the entire room - freezing all 4 players -
// the moment anyone clicked "Play Again" (RESTART_GAME) or a client sent
// START_GAME. The *Unlocked variants are what MapMessageToGameAction must call.
func TestBug3_RestartAndStartUnlocked_NoDeadlockUnderHeldMutex(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)
    p1, _ := gm.AddPlayer("P1")
    gm.AddPlayer("P2")
    gm.AddPlayer("P3")
    gm.AddPlayer("P4")

    done := make(chan error, 1)
    go func() {
        // Mirrors ws/client.go MapMessageToGameAction: lock once, then dispatch.
        gm.Mutex.Lock()
        defer gm.Mutex.Unlock()
        done <- gm.StartGameUnlocked(p1.ID)
    }()

    select {
    case err := <-done:
        if err != nil {
            t.Fatalf("StartGameUnlocked failed: %v", err)
        }
    case <-time.After(2 * time.Second):
        t.Fatal("Bug 3 STILL BROKEN: StartGameUnlocked deadlocked while gm.Mutex was already held")
    }

    // Force game into FINISHED so RestartGame's precondition is met.
    gm.Game.Status = models.StateFinished

    done2 := make(chan error, 1)
    go func() {
        gm.Mutex.Lock()
        defer gm.Mutex.Unlock()
        done2 <- gm.RestartGameUnlocked(p1.ID)
    }()

    select {
    case err := <-done2:
        if err != nil {
            t.Fatalf("RestartGameUnlocked failed: %v", err)
        }
    case <-time.After(2 * time.Second):
        t.Fatal("Bug 3 STILL BROKEN: RestartGameUnlocked deadlocked while gm.Mutex was already held")
    }

    t.Log("Bug 3 FIXED: StartGameUnlocked/RestartGameUnlocked run safely under an already-held mutex")
}

// Bug 4: a player's IsConnected flag never left true once set in AddPlayer,
// even after their WebSocket connection dropped. The frontend renders this
// flag as a green/gray status dot for every seat (Table.tsx), so the other
// 3 players had no way to tell a teammate had disconnected - the dot stayed
// green forever. ws/client.go's ReadPump now calls SetPlayerConnectedUnlocked
// on disconnect; this test guards the GameManager-level piece of that fix.
func TestBug4_SetPlayerConnected_ReflectsDisconnect(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)
    p1, _ := gm.AddPlayer("P1")

    if !gm.Game.Players[0].IsConnected {
        t.Fatalf("player should start connected")
    }

    if !gm.SetPlayerConnected(p1.ID, false) {
        t.Fatalf("Bug 4 STILL BROKEN: SetPlayerConnected did not find known player %s", p1.ID)
    }
    if gm.Game.Players[0].IsConnected {
        t.Fatalf("Bug 4 STILL BROKEN: IsConnected still true after disconnect")
    }

    if gm.SetPlayerConnected("no-such-player", false) {
        t.Fatalf("SetPlayerConnected should report false for an unknown player ID")
    }

    t.Log("Bug 4 FIXED: SetPlayerConnected flips IsConnected so other players' status dot updates")
}

// Bug 4b: ReadPump's disconnect handler calls SetPlayerConnectedUnlocked from
// inside a block that already holds gm.Mutex (to keep the broadcast read
// race-free against concurrent, already-locked MapMessageToGameAction
// mutations). Guard against that pattern deadlocking, mirroring Bug 3's test.
func TestBug4b_SetPlayerConnectedUnlocked_NoDeadlockUnderHeldMutex(t *testing.T) {
    store := &MockStore{}
    gm := NewGame(store)
    p1, _ := gm.AddPlayer("P1")

    done := make(chan bool, 1)
    go func() {
        gm.Mutex.Lock()
        defer gm.Mutex.Unlock()
        done <- gm.SetPlayerConnectedUnlocked(p1.ID, false)
    }()

    select {
    case found := <-done:
        if !found {
            t.Fatalf("SetPlayerConnectedUnlocked did not find known player %s", p1.ID)
        }
    case <-time.After(2 * time.Second):
        t.Fatal("Bug 4b STILL BROKEN: SetPlayerConnectedUnlocked deadlocked while gm.Mutex was already held")
    }
}
