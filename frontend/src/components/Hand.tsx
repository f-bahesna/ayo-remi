import React, { useState, useEffect } from 'react';
import Card from './Card';
import type { Card as CardType } from '../types';
import type { TurnPhase } from '../types';
import './Hand.css';

interface HandProps {
    cards: CardType[];
    onPlaySet: (cards: CardType[]) => void;
    onDiscard: (cardId: string) => void;
    isMyTurn: boolean;
    turnPhase: TurnPhase;
    onSelectionChange?: (cardIds: string[]) => void;
}

const Hand: React.FC<HandProps> = ({ cards, onPlaySet, onDiscard, isMyTurn, turnPhase, onSelectionChange }) => {
    const [selectedIndices, setSelectedIndices] = useState<number[]>([]);

    // Clear selection when hand changes (cards added/removed) or phase transitions
    useEffect(() => {
        setSelectedIndices([]);
    }, [cards.length, turnPhase]);

    // Notify parent of selection changes
    useEffect(() => {
        if (onSelectionChange) {
            const selectedIds = selectedIndices.map(i => cards[i]?.id).filter(Boolean);
            onSelectionChange(selectedIds);
        }
    }, [selectedIndices, cards]);

    const toggleSelect = (index: number) => {
        let next: number[];
        if (selectedIndices.includes(index)) {
            next = selectedIndices.filter(i => i !== index);
        } else {
            next = [...selectedIndices, index];
        }
        setSelectedIndices(next);
    };

    const handlePlaySet = () => {
        const selectedCards = selectedIndices.map(i => cards[i]);
        onPlaySet(selectedCards);
        setSelectedIndices([]);
    };

    const handleDiscard = () => {
        if (selectedIndices.length !== 1) return;
        const card = cards[selectedIndices[0]];
        onDiscard(card.id);
        setSelectedIndices([]);
    }

    return (
        <div className="hand-container">
            <div className="hand-controls">
                <button 
                    className="action-button primary" 
                    onClick={handlePlaySet} 
                    disabled={!isMyTurn || turnPhase === 'DRAW' || selectedIndices.length < 3}
                >
                    Play Set
                </button>
                <button 
                    className="action-button discard-btn" 
                    onClick={handleDiscard} 
                    disabled={!isMyTurn || turnPhase === 'DRAW' || selectedIndices.length !== 1}
                >
                    Discard (End Turn)
                </button>
            </div>
            <div className="hand-cards">
                {cards.map((card, index) => (
                    <div key={card.id || index} className="hand-card-wrapper" style={{ '--index': index } as React.CSSProperties}>
                         <Card 
                            card={card} 
                            selected={selectedIndices.includes(index)}
                            onClick={() => toggleSelect(index)}
                            draggable={isMyTurn && turnPhase !== 'DRAW'}
                            onDragStart={(e) => {
                                e.dataTransfer.setData("text/plain", card.id);
                                e.dataTransfer.effectAllowed = "move";
                            }}
                         />
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Hand;
