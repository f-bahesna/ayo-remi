import React, { useEffect } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import Table from '../components/Table';
import { WebSocketProvider, useGame } from '../context/WebSocketContext';
import '../App.css';

const GameContent: React.FC = () => {
    const { connect, isConnected } = useGame();
    const [searchParams] = useSearchParams();
    const { roomId } = useParams();
    const name = searchParams.get('name');

    useEffect(() => {
        if (name && roomId && !isConnected) {
            connect(name, roomId);
        }
    }, [name, roomId, isConnected, connect]);

    if (!name || !roomId) {
        return <div className="loading">Missing Name or Room ID</div>;
    }

    return (
        <div className="app">
            <h1 className="game-title">Remi Room: {roomId.slice(0, 8)}...</h1>
             {/* Copy Link Helper */}
             <div style={{ textAlign: 'center', marginBottom: '10px' }}>
                <button 
                    className="action-button secondary" 
                    style={{ fontSize: '0.8rem', padding: '5px 10px' }}
                    onClick={() => {
                        const url = window.location.href.split('?')[0]; // Share base URL to room
                        navigator.clipboard.writeText(url);
                        alert("Room Link Copied! Share it with friends.");
                    }}
                >
                    Copy Room Link
                </button>
            </div>
            
            <Table roomId={roomId || ''} />
        </div>
    );
};

const Room: React.FC = () => {
    return (
        <WebSocketProvider>
            <GameContent />
        </WebSocketProvider>
    );
};

export default Room;
