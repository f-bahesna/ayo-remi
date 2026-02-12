import React from 'react';
import { useGame } from '../context/WebSocketContext';
import Card from './Card';
import Hand from './Hand';
import './Table.css';

interface TableProps {
    roomId: string;
}

const Table: React.FC<TableProps> = ({ roomId }) => {
    const { gameState, drawCard, drawFromPile, playSet, discardCard, declareWin, startGame, restartGame, isMyTurn } = useGame();
    const [selectedHandCardIds, setSelectedHandCardIds] = React.useState<string[]>([]);

    if (!gameState) return <div className="loading">Loading game state...</div>;

    const isRoomMaster = gameState.mySeatIndex === 0;

    // ... existing handlers ...
    // To keep the file content intact, I'll assume handlers are here or unchanged
    // But I must include them if I am replacing the block containing them.
    // The previous view_file showed handlers defined inside component.
    // I will use specific replacement for the functional component signature/start and the waiting block.
    // Wait, replace_file_content needs contiguous block.
    // I can replace from "const Table..." down to the end of Waiting block.
    
    // Let's redefine handlers here to be safe, or just use "... handlers ..." if I had multi-replace.
    // But I must provide exact content.
    // I will replace the start of the component and the Waiting block.
    
    const handleDrawDeck = () => {
        if (isMyTurn && gameState.deckCount > 0 && gameState.turnPhase === 'DRAW') {
            drawCard("DECK");
        }
    }



    if (gameState.status === 'FINISHED') {
         // ... (Keep existing FINISHED block code)
         return (
            <div className="game-table">
                <div className="center-area">
                    <div className="waiting-overlay">
                         <h2 style={{ fontSize: '3rem', color: '#fbbf24' }}>GAME OVER</h2>
                         {gameState.winnerId ? (
                            <h3 style={{ fontSize: '2rem' }}>Winner: {gameState.winnerId === gameState.mySeatIndex.toString() ? "YOU!" : `Player ${gameState.players?.find(p => p.id === gameState.winnerId)?.name || gameState.winnerId.slice(0,8)}`}</h3>
                         ) : <h2>Winner Declared!</h2>}
                         
                         <div className="scoreboard" style={{ margin: '2rem 0', textAlign: 'left' }}>
                            <h3>Scores:</h3>
                            {gameState.players?.map(p => (
                                <div key={p.id} style={{ display: 'flex', justifyContent: 'space-between', width: '200px' }}>
                                    <span>{p.name}</span>
                                    <span>{p.score} pts</span>
                                </div>
                            ))}
                         </div>

                         {isRoomMaster && (
                             <button 
                                className="action-button primary" 
                                onClick={restartGame}
                                style={{ marginTop: '1rem', fontSize: '1.2rem', padding: '10px 20px' }}
                             >
                                 Play Again (Restart)
                             </button>
                         )}
                         {!isRoomMaster && <p>Waiting for Master to restart...</p>}
                    </div>
                </div>
            </div>
        );
    }

    // State for notifications
    const [notification, setNotification] = React.useState<string | null>(null);
    const prevPlayerCount = React.useRef(0);

    // Effect to detect new players
    React.useEffect(() => {
        if (gameState && gameState.players) {
            const currentCount = gameState.players.length;
            if (prevPlayerCount.current > 0 && currentCount > prevPlayerCount.current) {
                // New player joined
                const newPlayer = gameState.players[currentCount - 1]; // Assuming appended
                setNotification(`${newPlayer.name} joined the game!`);
                setTimeout(() => setNotification(null), 3000);
            }
            prevPlayerCount.current = currentCount;
        }
    }, [gameState?.players]);

    if (gameState.status === 'WAITING') {
        const slotsFilled = gameState.players ? gameState.players.length : (gameState.opponentHandSizes.filter(s => s >= 0).length + 1);
        const shareLink = `${window.location.protocol}//${window.location.host}/?room=${roomId}`;
        
        return (
            <div className="game-table">
                {notification && (
                    <div className="notification-toast" style={{
                        position: 'fixed', top: '20px', right: '20px', 
                        background: '#10b981', color: 'white', padding: '10px 20px', 
                        borderRadius: '8px', zIndex: 1000, boxShadow: '0 4px 6px rgba(0,0,0,0.2)',
                        animation: 'slideIn 0.3s ease-out'
                    }}>
                        {notification}
                    </div>
                )}
                
                <div className="center-area">
                    <div className="waiting-overlay" style={{ textAlign: 'center', background: 'rgba(0,0,0,0.8)', padding: '2rem', borderRadius: '1rem', backdropFilter: 'blur(10px)', border: '1px solid rgba(255,255,255,0.1)', maxWidth: '500px', width: '90%' }}>
                        <h2 style={{ fontSize: '2rem', marginBottom: '1rem' }}>Waiting for Players...</h2>
                        <div className="loading-spinner"></div>
                        
                        <div style={{ margin: '1.5rem 0', background: 'rgba(255,255,255,0.1)', padding: '1rem', borderRadius: '8px' }}>
                            <p style={{ fontSize: '0.9rem', color: '#94a3b8', marginBottom: '0.5rem' }}>Room ID:</p>
                            <code style={{ display: 'block', fontSize: '1.2rem', fontWeight: 'bold', color: '#fbbf24', marginBottom: '1rem' }}>{roomId}</code>
                            
                            <p style={{ fontSize: '0.9rem', color: '#94a3b8', marginBottom: '0.5rem' }}>Invite Link:</p>
                            <div style={{ display: 'flex', gap: '0.5rem' }}>
                                <input 
                                    readOnly 
                                    value={shareLink} 
                                    onClick={(e) => e.currentTarget.select()}
                                    style={{ background: '#1e293b', border: '1px solid #475569', color: 'white', padding: '0.5rem', borderRadius: '4px', flex: 1 }}
                                />
                                <button 
                                    className="action-button" 
                                    onClick={() => {
                                        if (navigator.clipboard && window.isSecureContext) {
                                            navigator.clipboard.writeText(shareLink)
                                                .then(() => {
                                                    setNotification("Link Copied!");
                                                    setTimeout(() => setNotification(null), 2000);
                                                })
                                                .catch(() => {
                                                    prompt("Copy this link:", shareLink);
                                                });
                                        } else {
                                            prompt("Copy this link:", shareLink);
                                        }
                                    }}
                                    title="Copy Link"
                                >
                                    Copy
                                </button>
                            </div>
                        </div>

                        <div className="joined-players-list" style={{ textAlign: 'left', background: 'rgba(0,0,0,0.2)', padding: '1rem', borderRadius: '8px', marginBottom: '1rem' }}>
                            <h4 style={{ margin: '0 0 0.5rem 0', color: '#94a3b8' }}>Players ({slotsFilled}/4):</h4>
                            <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                                {gameState.players?.map((p, i) => (
                                    <li key={p.id || i} style={{ padding: '0.5rem', borderBottom: '1px solid rgba(255,255,255,0.05)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                        <div style={{ width: '8px', height: '8px', borderRadius: '50%', background: p.isConnected ? '#10b981' : '#64748b' }}></div>
                                        <span>{p.name} {p.id === gameState.mySeatIndex.toString() ? "(You)" : ""}</span>
                                        {p.seatIndex === 0 && <span style={{ fontSize: '0.8rem', background: '#eab308', color: 'black', padding: '2px 6px', borderRadius: '4px', marginLeft: 'auto' }}>Master</span>}
                                    </li>
                                ))}
                            </ul>
                        </div>
                        
                         {isRoomMaster && (
                            <div style={{ marginTop: '2rem' }}>
                                <p style={{ color: '#fbbf24', fontWeight: 'bold', marginBottom: '0.5rem' }}>You are the Master</p>
                                <p style={{ fontSize: '0.9rem', color: '#94a3b8' }}>Waiting for {4 - slotsFilled} more player(s) to auto-start...</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        );
    }

    // Helper to map seat index to UI position (Bottom, Left, Top, Right) relative to 'mySeatIndex'
    const getPosition = (seatIndex: number) => {
        const diff = (seatIndex - gameState.mySeatIndex + 4) % 4;
        if (diff === 0) return 'bottom'; // Me
        if (diff === 1) return 'left';
        if (diff === 2) return 'top';
        if (diff === 3) return 'right';
        return 'unknown';
    };

    const renderOpponent = (seatIndex: number, handSize: number) => {
        const pos = getPosition(seatIndex);
        const isActive = gameState.currentTurnPlayer === seatIndex;
        // Find player data if available (e.g. name, score)
        // Note: gameState.players might not be fully populated in PublicView unless we add it.
        // PublicView usually has "OpponentHandSizes". Name/Score might be missing.
        // Let's assume PublicView strictly follows `state.go` struct.
        // Actually, `PublicGameView` does NOT have `Players` list, only `OpponentHandSizes`.
        // So we can't show names easily unless we add `Players` metadata to PublicView.
        // For now, use "Player X".
        
        const player = gameState.players?.find(p => p.seatIndex === seatIndex);

        return (
            <div className={`opponent-seat opponent-${pos} ${isActive ? 'active' : ''}`} key={seatIndex}>
                <div className="avatar">P{seatIndex + 1}</div>
                <div className="opponent-hand" data-count={handSize}>
                    {/* Render minimal card backs for opponent hand */}
                    {Array.from({ length: Math.min(handSize, 5) }).map((_, i) => (
                         <div key={i} className="mini-card-back"></div>
                    ))}
                    <span className="hand-count">{handSize}</span>
                </div>
                {/* Render Played Sets */}
                {player && player.playedSets && player.playedSets.length > 0 && (
                    <div className="player-played-sets">
                        {player.playedSets.map((set, setIdx) => (
                            <div key={setIdx} className="card-group small">
                                {set.map((card, cIdx) => (
                                    <Card key={cIdx} card={card} style={{ transform: 'scale(0.6)', margin: '-10px -15px' }} />
                                ))}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        );
    };

    // const topPileCard = gameState.pile.length > 0 ? gameState.pile[gameState.pile.length - 1] : null;

    return (
        <div className="game-table">
            {/* Center Area */}
            <div 
                className="center-area"
                onDragOver={(e) => {
                    e.preventDefault(); // Allow drop
                    e.dataTransfer.dropEffect = "move";
                }}
                onDrop={(e) => {
                    e.preventDefault();
                    if (!isMyTurn || gameState.turnPhase === 'DRAW') return;
                    const cardId = e.dataTransfer.getData("text/plain");
                    if (cardId) {
                        discardCard(cardId);
                    }
                }}
            >
                {/* Deck */}
                <div 
                    className={`deck-pile ${isMyTurn && gameState.turnPhase === 'DRAW' ? 'interactive' : ''}`} 
                    onClick={handleDrawDeck}
                >
                    {gameState.deckCount > 0 ? (
                        <div className="deck-card">
                             <div className="card-back-pattern"></div>
                        </div>
                    ) : (
                        <div className="empty-slot">Empty (Reshuffles)</div>
                    )}
                    <span className="deck-count">{gameState.deckCount} Cards</span>
                </div>

                {/* Pile */}
                {/* Pile (Radial Layout) */}
                 <div className={`discard-pile ${isMyTurn && gameState.turnPhase === 'DRAW' ? 'interactive' : ''}`}>
                    {gameState.pile.length > 0 ? (
                        gameState.pile.map((card, index) => {
                             
                             const angle = index * 25;
                             const radius = 140; 
                             
                             const simpleTransform = `translate(-50%, -50%) rotate(${angle}deg) translate(0, -${radius}px)`;

                             const isLast3 = (gameState.pile.length - index) <= 3;
                             const hasHandSelection = selectedHandCardIds.length >= 2;
                             const isEligible = isMyTurn && gameState.turnPhase === 'DRAW' && isLast3;

                            return (
                                <Card 
                                    key={card.id || index} 
                                    card={card} 
                                    onClick={() => {
                                        if (isMyTurn && gameState.turnPhase === 'DRAW') {
                                            if (hasHandSelection) {
                                                // Pick specific card from pile (must be last 3)
                                                if (!isLast3) {
                                                    alert("Hanya bisa mengambil dari 3 kartu terakhir di pile.");
                                                    return;
                                                }
                                                drawFromPile(selectedHandCardIds, card.id);
                                            } else {
                                                // Existing behavior: draw from top of pile
                                                const count = gameState.pile.length - index;
                                                if (count <= 3) {
                                                    if (count > 1) {
                                                         if (confirm(`Draw last ${count} cards? You must form a set with ALL of them immediately.`)) {
                                                             drawCard("PILE", count);
                                                         }
                                                    } else {
                                                        drawCard("PILE", 1);
                                                    }
                                                } else {
                                                    alert("You can only draw up to 3 cards from the pile.");
                                                }
                                            }
                                        }
                                    }}
                                    className={`pile-card ${hasHandSelection && isEligible ? 'eligible-pick' : ''}`}
                                    style={{ 
                                        cursor: isEligible ? 'pointer' : 'default',
                                        opacity: (isMyTurn && gameState.turnPhase === 'DRAW' && isLast3) ? 1 : 0.6,
                                        transform: simpleTransform,
                                        zIndex: index,
                                        '--base-transform': simpleTransform
                                    } as React.CSSProperties}
                                />
                             );
                        })
                    ) : null}
                </div>

                </div>


            {/* Players */}
            {gameState.opponentHandSizes.map((size, index) => {
                if (index === gameState.mySeatIndex) return null;
                return renderOpponent(index, size);
            })}

            {/* My Hand (Bottom) */}
            <div className={`my-seat ${isMyTurn ? 'active-turn' : ''} phase-${gameState.turnPhase?.toLowerCase()}`}>
                <div className="my-actions">
                     <span className="turn-indicator">
                        {isMyTurn 
                            ? `YOUR TURN - ${gameState.turnPhase === 'DRAW' ? 'Draw Card' : gameState.turnPhase === 'PLAY' ? 'Play Sets or Discard' : 'Discard to End'}`
                            : `Player ${gameState.currentTurnPlayer + 1}'s Turn`
                        }
                     </span>
                     
                     {/* Score Display (during game?) */}
                     {gameState.score !== undefined && <span className="score-display">Score: {gameState.score}</span>}

                </div>
                
                {/* Nutup Button - Bottom Right Absolute Position */}
                {isMyTurn && gameState.turnPhase !== 'DRAW' && (
                    <button 
                        className="action-button nutup" 
                        onClick={declareWin}
                        // Allow clicking to check win condition on backend.
                        // We could duplicate IsWinningHand logic here for better UX, but for now let backend validate.
                        // Removing the strict frontend checks.
                        title="Declare Win (Nutup)"
                        style={{
                            position: 'absolute',
                            bottom: '20px',
                            right: '20px',
                            zIndex: 100,
                            padding: '12px 24px',
                            fontSize: '1.2rem',
                            fontWeight: 'bold',
                            boxShadow: '0 4px 6px rgba(0,0,0,0.3)'
                        }}
                    >
                        Nutup (Win)
                    </button>
                 )}


                {/* My Played Sets */}
                {(() => {
                    const me = gameState.players?.find(p => p.seatIndex === gameState.mySeatIndex);
                    if (me && me.playedSets && me.playedSets.length > 0) {
                        return (
                            <div className="player-played-sets my-sets">
                                {me.playedSets.map((set, setIdx) => (
                                    <div key={setIdx} className="card-group small">
                                        {set.map((card, cIdx) => (
                                            <Card key={cIdx} card={card} style={{ transform: 'scale(0.7)', margin: '-8px -12px' }} />
                                        ))}
                                    </div>
                                ))}
                            </div>
                        );
                    }
                    return null;
                })()}
                
                <Hand 
                    cards={gameState.myHand || []} 
                    onPlaySet={playSet} 
                    onDiscard={(id) => {
                        discardCard(id);
                    }}
                    isMyTurn={isMyTurn}
                    turnPhase={gameState.turnPhase}
                    onSelectionChange={setSelectedHandCardIds}
                />
            </div>
        </div>
    );
};

export default Table;
