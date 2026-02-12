import React from 'react';
import type { Card as CardType } from '../types';
import './Card.css';

interface CardProps {
    card?: CardType; // If undefined, it's face down
    onClick?: () => void;
    selected?: boolean;
    style?: React.CSSProperties;
    className?: string; // Allow passing extra classes
    draggable?: boolean;
    onDragStart?: (e: React.DragEvent) => void;
}

const Card: React.FC<CardProps> = ({ card, onClick, selected, style, className, draggable, onDragStart }) => {
    const isFaceUp = !!card;

    return (
        <div 
            className={`card ${isFaceUp ? 'face-up' : 'face-down'} ${selected ? 'selected' : ''} ${className || ''}`} 
            onClick={onClick}
            style={style}
            draggable={draggable}
            onDragStart={onDragStart}
        >
            {isFaceUp ? (
                <div className="card-content" data-suit={card.suit}>
                    <div className="card-top-left">
                        <span>{getRankLabel(card.rank)}</span>
                        <span className="suit">{getSuitIcon(card.suit)}</span>
                    </div>
                    <div className="card-center">
                        <span className="suit-large">{getSuitIcon(card.suit)}</span>
                    </div>
                    <div className="card-bottom-right">
                        <span>{getRankLabel(card.rank)}</span>
                        <span className="suit">{getSuitIcon(card.suit)}</span>
                    </div>
                </div>
            ) : (
                <div className="card-back">
                    <div className="pattern"></div>
                </div>
            )}
        </div>
    );
};

// Helpers
function getRankLabel(rank: number): string {
    if (rank === 11) return 'J';
    if (rank === 12) return 'Q';
    if (rank === 13) return 'K';
    if (rank === 14) return 'A';
    if (rank === 1) return 'A'; // Treated as 1 or A depending on logic, verify with user if 1 shows '1' or 'A'. Usually 'A'.
    if (rank === 0) return 'JK';
    return rank.toString();
}

function getSuitIcon(suit: string): string {
    switch (suit) {
        case 'S': return '♠';
        case 'H': return '♥';
        case 'D': return '♦';
        case 'C': return '♣';
        case 'J': return '★'; // Joker
        default: return '?';
    }
}

export default Card;
