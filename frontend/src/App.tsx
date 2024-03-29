import React from 'react';
import { Content } from './components/Content/Content';
import { Topbar } from './components/Topbar/Topbar';
import { GameProvider } from './context/GameContext';
import { UserProvider } from './context/UserContext';
import OnlineMultiplayerGameService from './service/OnlineMultiplayerGameService';

export const App: React.FC = () => {
    return <>
        <UserProvider>
            <Topbar />
            <GameProvider gameServiceProvider={OnlineMultiplayerGameService}>
                <Content />
            </GameProvider>
        </UserProvider>
    </>
}
