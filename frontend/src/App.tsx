import React, { useState } from 'react';
import { useGame } from './context/WebSocketContext';
import Table from './components/Table';
import './App.css';

const AppContent: React.FC = () => {
    const { isConnected, connect, gameState } = useGame();
    const [name, setName] = useState('Player ' + Math.floor(Math.random() * 1000));
    const [roomId, setRoomId] = useState(() => {
        const params = new URLSearchParams(window.location.search);
        return params.get('room') || '';
    });

    const createRoom = async () => {
        try {
            const apiProtocol = window.location.protocol;
            const apiHost = window.location.hostname;
            const res = await fetch(`${apiProtocol}//${apiHost}:8080/api/rooms`, { method: 'POST' });
            if (!res.ok) throw new Error("Failed to create room");
            const data = await res.json();
            setRoomId(data.roomId);
            
            // Clean URL? Or update it?
            // Let's update URL so user sees it change
            const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname + '?room=' + data.roomId;
            window.history.pushState({path:newUrl},'',newUrl);
            
            connect(name, data.roomId);
        } catch (e) {
            console.error(e);
            alert('Error creating room');
        }
    };

    const joinRoom = () => {
        if (!roomId) return alert("Enter Room ID");
        connect(name, roomId);
    };

    if (!isConnected) {
        return (
            <div className="login-screen">
                <div className="login-card">
                    <h1>Remi Card Game</h1>
                    <div className="input-group">
                        <label>Enter Name</label>
                        <input 
                            value={name} 
                            onChange={e => setName(e.target.value)} 
                            placeholder="Your Name"
                        />
                    </div>
                    
                    <div className="input-group">
                        <label>Room ID</label>
                        <input 
                            value={roomId} 
                            onChange={e => setRoomId(e.target.value)} 
                            placeholder="Room ID (or click Create)"
                        />
                    </div>

                    <div className="button-group">
                        <button className="start-btn" onClick={createRoom}>
                            Create Room
                        </button>
                        <button className="start-btn secondary" onClick={joinRoom}>
                            Join Room
                        </button>
                    </div>
                    <p className="hint">Game starts when 4 players join.</p>
                </div>
            </div>
        );
    }

    if (!gameState) {
        return <div className="loading-screen">Connecting to game...</div>;
    }

    return <Table roomId={roomId} />;
};

import { WebSocketProvider } from './context/WebSocketContext';

// ... AppContent definition ...

function App() {
  return (
    <WebSocketProvider>
      <AppContent />
    </WebSocketProvider>
  );
}

export default App;
