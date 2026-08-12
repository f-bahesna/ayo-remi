import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { WebSocketProvider, useGame } from './WebSocketContext';

// Minimal fake WebSocket so connect() can run without a real network/server.
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

describe('WebSocketContext drawCard', () => {
    beforeEach(() => {
        FakeWebSocket.instances = [];
        // @ts-expect-error - test stub, not a full WebSocket implementation
        global.WebSocket = FakeWebSocket;
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    it('sends the requested count when drawing multiple cards from the pile', () => {
        const { result } = renderHook(() => useGame(), {
            wrapper: WebSocketProvider,
        });

        act(() => {
            result.current.connect('Alice', 'room1');
        });

        const socket = FakeWebSocket.instances[0];
        act(() => {
            socket.onopen?.();
        });

        act(() => {
            result.current.drawCard('PILE', 3);
        });

        // The last sent message should be the DRAW_CARD action (JOIN_GAME is sent first on open).
        const lastMessage = JSON.parse(socket.sent[socket.sent.length - 1]);
        expect(lastMessage.type).toBe('DRAW_CARD');
        expect(lastMessage.payload).toEqual({ source: 'PILE', count: 3 });
    });
});
