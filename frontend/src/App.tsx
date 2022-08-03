import React from 'react';
import { Content } from './components/Content/Content';
import { Topbar } from './components/Topbar/Topbar';
import { GameProvider } from './context/GameContext';
import OfflineMultiplayerGameService from './service/OfflineMultiplayerGameService';

export const App: React.FC = () => {
    return <>
        <Topbar/>
        <GameProvider gameService={OfflineMultiplayerGameService}>
            <Content/>
        </GameProvider>
    </>
}