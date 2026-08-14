import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, screen, act } from '@testing-library/react';
import { WebSocketProvider } from '../context/WebSocketContext';
import { useGame } from '../context/useGame';
import Table from './Table';
import type { GameState } from '../types';

// Minimal fake WebSocket so connect() can run without a real network/server,
// mirroring the harness in WebSocketContext.test.tsx.
class FakeWebSocket {
    static instances: FakeWebSocket[] = [];
    static OPEN = 1;

    url: string;
    readyState = FakeWebSocket.OPEN;
    sent: string[] = [];
    onopen: (() => void) | null = null;
    onmessage: ((event: { data: string }) => void) | null = null;
    onclose: (() => void) | null = null;

    constructor(url: string) {
        this.url = url;
        FakeWebSocket.instances.push(this);
    }

    send(data: string) {
        this.sent.push(data);
    }

    close() {
        this.onclose?.();
    }
}

// Drives the provider's connect() then delivers a GAME_UPDATE payload, the same
// way the real backend pushes state after every action.
const Harness: React.FC<{ state: GameState }> = ({ state }) => {
    const { connect, gameState } = useGame();
    React.useEffect(() => {
        connect('Alice', 'room1');
        const socket = FakeWebSocket.instances[FakeWebSocket.instances.length - 1];
        socket.onopen?.();
        socket.onmessage?.({ data: JSON.stringify({ type: 'GAME_UPDATE', payload: state }) });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    if (!gameState) return null;
    return <Table roomId="room1" />;
};

import React from 'react';

const basePlayers = [
    { id: 'p0-uuid', name: 'Alice', seatIndex: 0, score: 0, isConnected: true, hasPlayedSet: false },
    { id: 'p1-uuid', name: 'Bob', seatIndex: 1, score: 0, isConnected: true, hasPlayedSet: false },
    { id: 'p2-uuid', name: 'Carol', seatIndex: 2, score: 0, isConnected: true, hasPlayedSet: false },
    { id: 'p3-uuid', name: 'Dave', seatIndex: 3, score: 0, isConnected: true, hasPlayedSet: false },
];

describe('Table player-identity display', () => {
    beforeEach(() => {
        FakeWebSocket.instances = [];
        // @ts-expect-error - test stub, not a full WebSocket implementation
        global.WebSocket = FakeWebSocket;
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it('marks the local player as "(You)" in the waiting room using seatIndex, not id vs. seatIndex string', () => {
        const state: GameState = {
            id: 'g1',
            status: 'WAITING',
            currentTurnPlayer: 0,
            turnPhase: 'DRAW',
            deckCount: 0,
            pile: [],
            tableSets: [],
            myHand: [],
            mySeatIndex: 0, // Alice is seat 0 locally; her id is 'p0-uuid', never "0"
            opponentHandSizes: [0, 0, 0, 0],
            hasTakenFromPile: false,
            masterPlayerId: 'p0-uuid',
            players: basePlayers,
            hasPlayedSet: false,
            score: 0,
        };

        render(
            <WebSocketProvider>
                <Harness state={state} />
            </WebSocketProvider>
        );

        act(() => {});

        const aliceRow = screen.getByText((content) => content.includes('Alice') && content.includes('(You)'));
        expect(aliceRow).toBeTruthy();
        expect(screen.queryByText((content) => content.includes('Bob') && content.includes('(You)'))).toBeNull();
    });

    it('shows "YOU!" as the winner when the winning player is the local seat, comparing by seatIndex not id', () => {
        const state: GameState = {
            id: 'g1',
            status: 'FINISHED',
            currentTurnPlayer: 0,
            turnPhase: 'PLAY',
            deckCount: 0,
            pile: [],
            tableSets: [],
            myHand: [],
            mySeatIndex: 0, // local player is Alice (id 'p0-uuid', seatIndex 0)
            opponentHandSizes: [0, 0, 0, 0],
            winnerId: 'p0-uuid',
            hasTakenFromPile: false,
            masterPlayerId: 'p0-uuid',
            players: basePlayers,
            hasPlayedSet: false,
            score: 0,
        };

        render(
            <WebSocketProvider>
                <Harness state={state} />
            </WebSocketProvider>
        );

        act(() => {});

        expect(screen.getByText(/YOU!/)).toBeTruthy();
    });
});
