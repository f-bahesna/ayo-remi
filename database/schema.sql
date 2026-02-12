-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status TEXT NOT NULL, -- 'WAITING', 'IN_PROGRESS', 'FINISHED'
    current_turn_player_index INT DEFAULT 0,
    deck JSONB, -- Storing remaining cards as JSON array
    pile JSONB DEFAULT '[]', -- Storing center pile as JSON array
    table_sets JSONB DEFAULT '[]', -- Public sets on table
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    seat_index INT NOT NULL,
    name TEXT NOT NULL,
    hand JSONB DEFAULT '[]', -- Private cards
    is_connected BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(game_id, seat_index)
);

CREATE INDEX idx_players_game_id ON players(game_id);
