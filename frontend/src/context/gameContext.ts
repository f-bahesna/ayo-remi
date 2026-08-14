import { createContext } from 'react';
import type { Card, GameState } from '../types';

export interface WebSocketContextType {
    isConnected: boolean;
    gameState: GameState | null;
    connect: (name: string, roomId?: string) => void;
    drawCard: (source?: "DECK" | "PILE", count?: number) => void;
    drawFromPile: (handCardIds: string[], pileCardId: string) => void;
    playSet: (cards: Card[]) => void;
    discardCard: (cardId: string) => void;
    declareWin: () => void;
    startGame: () => void;
    restartGame: () => void;
    isMyTurn: boolean;
}

export const WebSocketContext = createContext<WebSocketContextType | null>(null);
