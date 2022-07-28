import React from 'react';
import { Content } from './Content/Content';
import "./Frontpage.css"
import { Topbar } from './Topbar/Topbar';

export const Frontpage: React.FC = () => {
    return <>
        <Topbar/>
        <Content/>
    </>
}