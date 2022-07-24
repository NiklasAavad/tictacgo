import React from 'react';
import "./Frontpage.css"

export const Frontpage: React.FC = () => {
    return <div className="container">
        <div className='top'>
            <div className='header'>Header</div>
            <div className='sub-header'>Lobby | User | Help</div>
        </div>
        <div className='center'>
            <div className='chat'>Chat</div>
            <div className='lobby'>Lobby</div>
        </div>
    </div>
}