import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import type { GameState, WSMessage } from '../types';

interface WebSocketContextType {
    isConnected: boolean;
    gameState: GameState | null;
    connect: (name: string, roomId?: string) => void;
    drawCard: (source?: "DECK" | "PILE", count?: number) => void;
    drawFromPile: (handCardIds: string[], pileCardId: string) => void;
    playSet: (cards: any[]) => void;
    discardCard: (cardId: string) => void;
    declareWin: () => void;
    startGame: () => void;
    restartGame: () => void;
    isMyTurn: boolean;
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [isConnected, setIsConnected] = useState(false);
    const [gameState, setGameState] = useState<GameState | null>(null);
    const ws = useRef<WebSocket | null>(null);

    useEffect(() => {
        return () => {
            if (ws.current) {
                ws.current.close();
            }
        };
    }, []);

    const connect = (name: string, roomId: string = 'default') => {
        if (ws.current) {
            ws.current.close();
        }

        // Use current hostname (e.g., actual IP) to connect to backend on port 8080
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsHost = window.location.hostname; 
        const wsUrl = `${wsProtocol}//${wsHost}:8080/ws?room=${roomId}&name=${name}`;
        console.log("Attempting WebSocket Connection to:", wsUrl);
        const socket = new WebSocket(wsUrl);
        ws.current = socket;

        socket.onopen = () => {
            console.log('Connected to WS');
            setIsConnected(true);
            // Send Join message
            socket.send(JSON.stringify({
                type: 'JOIN_GAME',
                payload: { name }
            }));
        };

        socket.onmessage = (event) => {
            try {
                const msg: WSMessage = JSON.parse(event.data);
                if (msg.type === 'GAME_UPDATE') {
                    setGameState(msg.payload);
                } else if (msg.type === 'ERROR') {
                    console.error('Game Error:', msg.payload);
                    alert(msg.payload.message); // Simple alert for now
                }
            } catch (e) {
                console.error('Failed to parse WS message', e);
            }
        };

        socket.onclose = () => {
            console.log('WS Disconnected');
            setIsConnected(false);
        };
    };

    const sendMessage = (type: string, payload: any = {}) => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify({ type, payload }));
        }
    };

    const drawCard = (source: "DECK" | "PILE" = "DECK") => sendMessage('DRAW_CARD', { source });
    const drawFromPile = (handCardIds: string[], pileCardId: string) => sendMessage('DRAW_FROM_PILE', { handCardIds, pileCardId });
    const playSet = (cards: any[]) => sendMessage('PLAY_SET', { cards });
    const discardCard = (cardId: string) => sendMessage('DISCARD_CARD', { cardId });
    const declareWin = () => sendMessage('DECLARE_WIN');
    const startGame = () => sendMessage('START_GAME');
    const restartGame = () => sendMessage('RESTART_GAME');

    const isMyTurn = gameState?.currentTurnPlayer === gameState?.mySeatIndex;

    return (
        <WebSocketContext.Provider value={{ isConnected, gameState, connect, drawCard, drawFromPile, playSet, discardCard, declareWin, startGame, restartGame, isMyTurn }}>
            {children}
        </WebSocketContext.Provider>
    );
};

export const useGame = () => {
    const context = useContext(WebSocketContext);
    if (!context) throw new Error('useGame must be used within WebSocketProvider');
    return context;
};
