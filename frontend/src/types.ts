export type Suit = "S" | "H" | "D" | "C" | "J";
export type Rank = number; // 1-9, 11=J, 12=Q, 13=K, 14=A, 0=Joker

export interface Card {
    suit: Suit;
    rank: Rank;
    id: string;
}

export interface Player {
    id: string;
    name: string;
    seatIndex: number;
    hand: Card[]; // Private
    hasTakenFromPile: boolean;
}

export interface PublicPlayer {
    id: string;
    name: string;
    seatIndex: number;
    score: number;
    isConnected: boolean;
    hasPlayedSet: boolean;
    playedSets?: Card[][];
}

export type GameStatus = "WAITING" | "IN_PROGRESS" | "FINISHED";
export type TurnPhase = "DRAW" | "PLAY" | "DISCARD";

export interface GameState {
    id: string;
    status: GameStatus;
    currentTurnPlayer: number; // Seat index 0-3
    turnPhase: TurnPhase;
    deckCount: number;
    pile: Card[];
    tableSets: Card[][];
    myHand: Card[];
    mySeatIndex: number;
    opponentHandSizes: number[]; // Ordered by seat index 0-3
    winnerId?: string;
    hasTakenFromPile: boolean;
    masterPlayerId: string;
    
    // Extra fields
    players: PublicPlayer[];
    hasPlayedSet: boolean;
    score: number;
}

export interface WSMessage {
    type: string;
    payload: any;
}
