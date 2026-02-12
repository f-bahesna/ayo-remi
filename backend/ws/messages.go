package ws

import "remi-game/models"

const (
    // Client -> Server
    MsgDrawCard     = "DRAW_CARD"
    MsgDrawFromPile = "DRAW_FROM_PILE"
    MsgPlaySet      = "PLAY_SET"
    MsgDiscardCard  = "DISCARD_CARD"
    MsgDeclareWin   = "DECLARE_WIN"
    MsgJoinGame     = "JOIN_GAME" 
    MsgStartGame    = "START_GAME"
    MsgRestartGame  = "RESTART_GAME"

    // Server -> Client
    MsgGameUpdate = "GAME_UPDATE"
    MsgError      = "ERROR"
)

type WSMessage struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

type DrawCardPayload struct {
    Source string `json:"source"` // "DECK" or "PILE"
    Count  int    `json:"count"`  // Default 1 if 0/missing
}

type DrawFromPilePayload struct {
    HandCardIDs []string `json:"handCardIds"` // IDs of cards selected from hand
    PileCardID  string   `json:"pileCardId"`  // ID of the card to pick from pile
}

type PlaySetPayload struct {
    Cards []models.Card `json:"cards"`
}

type DiscardCardPayload struct {
    CardID string `json:"cardId"`
}

type JoinGamePayload struct {
    Name string `json:"name"`
}

type StartGamePayload struct {
}

type ErrorPayload struct {
    Message string `json:"message"`
}
