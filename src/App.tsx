import React from 'react';
import { Content } from './components/Content/Content';
import { Topbar } from './components/Topbar/Topbar';

export const App: React.FC = () => {
    return <>
        <Topbar/>
        <Content/>
    </>
}