import React from 'react';
import { Content } from './components/Content/Content';
import { Topbar } from './components/Topbar/Topbar';
import { GameProvider } from './context/GameContext';
import { UserProvider } from './context/UserContext';
import OfflineMultiplayerGameService from './service/OfflineMultiplayerGameService';
import OfflineSingleplayerGameService from './service/OfflineSingleplayerGameService';

export const App: React.FC = () => {
    return <>
        <UserProvider>
            <Topbar />
            <GameProvider gameServiceProvider={OfflineSingleplayerGameService}>
                <Content />
            </GameProvider>
        </UserProvider>
    </>
}