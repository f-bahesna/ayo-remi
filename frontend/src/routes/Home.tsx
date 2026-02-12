import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../App.css'; // Reusing App styles

const Home: React.FC = () => {
    const [name, setName] = useState('');
    const [roomId, setRoomId] = useState('');
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    const handleCreateRoom = async () => {
        if (!name) {
            alert("Please enter your name");
            return;
        }
        setLoading(true);
        try {
            const res = await fetch('http://localhost:8080/api/rooms', {
                method: 'POST'
            });
            const data = await res.json();
            if (data.roomId) {
                navigate(`/room/${data.roomId}?name=${encodeURIComponent(name)}`);
            }
        } catch (err) {
            console.error(err);
            alert("Failed to create room");
        } finally {
            setLoading(false);
        }
    };

    const handleJoinRoom = () => {
        if (!name || !roomId) {
            alert("Please enter name and room ID");
            return;
        }
        navigate(`/room/${roomId}?name=${encodeURIComponent(name)}`);
    };

    return (
        <div className="game-container">
            <h1 className="game-title">Remi Online</h1>
            <div className="login-card">
                <input 
                    type="text" 
                    placeholder="Enter Your Name" 
                    value={name} 
                    onChange={(e) => setName(e.target.value)}
                    className="login-input"
                />
                
                <div className="divider">
                    <span>CREATE A ROOM</span>
                </div>
                
                <button 
                    className="login-button primary" 
                    onClick={handleCreateRoom}
                    disabled={loading}
                >
                    {loading ? "Creating..." : "Create New Room"}
                </button>

                <div className="divider">
                    <span>OR JOIN EXISTING</span>
                </div>

                <input 
                    type="text" 
                    placeholder="Enter Room ID" 
                    value={roomId} 
                    onChange={(e) => setRoomId(e.target.value)}
                    className="login-input"
                />
                <button className="login-button secondary" onClick={handleJoinRoom}>
                    Join Room
                </button>
            </div>
        </div>
    );
};

export default Home;
