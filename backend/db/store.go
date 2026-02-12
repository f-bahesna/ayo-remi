package db

import (
	"database/sql"
	"encoding/json"
	"remi-game/models"

	_ "github.com/lib/pq"
)

type GameStore struct {
	DB *sql.DB
}

func NewGameStore(connStr string) (*GameStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
        return nil, err
    }
	return &GameStore{DB: db}, nil
}

func (s *GameStore) SaveGame(game models.Game) error {
	deckJSON, _ := json.Marshal(game.Deck)
	pileJSON, _ := json.Marshal(game.Pile)
	tableSetsJSON, _ := json.Marshal(game.TableSets)

	// Upsert Game
    // Note: ON CONFLICT needs ID. Assuming ID is set.
	query := `
		INSERT INTO games (id, status, current_turn_player_index, deck, pile, table_sets, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			current_turn_player_index = EXCLUDED.current_turn_player_index,
			deck = EXCLUDED.deck,
			pile = EXCLUDED.pile,
			table_sets = EXCLUDED.table_sets,
			updated_at = NOW();
	`
	_, err := s.DB.Exec(query, game.ID, game.Status, game.CurrentTurnPlayer, deckJSON, pileJSON, tableSetsJSON)
	if err != nil {
		return err
	}

	// Upsert Players
	for _, p := range game.Players {
		handJSON, _ := json.Marshal(p.Hand)
		pQuery := `
			INSERT INTO players (id, game_id, seat_index, name, hand, is_connected, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			ON CONFLICT (id) DO UPDATE SET
				hand = EXCLUDED.hand,
				is_connected = EXCLUDED.is_connected,
				updated_at = NOW();
		`
		_, err := s.DB.Exec(pQuery, p.ID, game.ID, p.SeatIndex, p.Name, handJSON, p.IsConnected)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GameStore) LoadGame(gameID string) (*models.Game, error) {
	var g models.Game
	var deckJSON, pileJSON, tableSetsJSON []byte

	query := `SELECT id, status, current_turn_player_index, deck, pile, table_sets FROM games WHERE id = $1`
	row := s.DB.QueryRow(query, gameID)
	err := row.Scan(&g.ID, &g.Status, &g.CurrentTurnPlayer, &deckJSON, &pileJSON, &tableSetsJSON)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(deckJSON, &g.Deck)
	json.Unmarshal(pileJSON, &g.Pile)
	json.Unmarshal(tableSetsJSON, &g.TableSets)

	// Load Players
	pQuery := `SELECT id, seat_index, name, hand, is_connected FROM players WHERE game_id = $1 ORDER BY seat_index`
	rows, err := s.DB.Query(pQuery, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Player
		var handJSON []byte
		if err := rows.Scan(&p.ID, &p.SeatIndex, &p.Name, &handJSON, &p.IsConnected); err != nil {
			return nil, err
		}
		json.Unmarshal(handJSON, &p.Hand)
		g.Players = append(g.Players, p)
	}

	return &g, nil
}
