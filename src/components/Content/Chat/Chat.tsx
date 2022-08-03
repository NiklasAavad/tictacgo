import React, { useEffect, useMemo, useState } from "react";
import { useGameContext } from "../../../context/GameContext"
import "./Chat.css"

enum ChatType {
    GAME_INFO,
    USER_MESSAGE
}

type ChatMessage = {
    text: string,
    type: ChatType
}

export const Chat: React.FC = () => {
    const [messages, setMessages] = useState<ChatMessage[]>([])
    const { latestGameInfoMessage: latestGameMessage } = useGameContext();

    useEffect(() => {
        const gameInfoMessage = {text: latestGameMessage, type: ChatType.GAME_INFO};
        setMessages([...messages, gameInfoMessage]);
    }, [latestGameMessage])

    const styledMessages = useMemo(() => {
        return messages.map((message, idx) => {
            if (message.type == ChatType.GAME_INFO) {
                return <div key={idx} className="chat-message">
                    <span className="game-info-message">Game Info: </span>
                    {message.text}
                </div>
            }
        });
    }, [messages])

    return <div className='chat'>
        <div className="chat-header underline">Chat</div>
        {styledMessages}
    </div>
};