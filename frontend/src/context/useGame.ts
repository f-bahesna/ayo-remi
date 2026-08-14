import { useContext } from 'react';
import { WebSocketContext } from './gameContext';

export const useGame = () => {
    const context = useContext(WebSocketContext);
    if (!context) throw new Error('useGame must be used within WebSocketProvider');
    return context;
};
