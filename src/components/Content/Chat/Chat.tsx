import React, { useEffect, useMemo, useState } from "react";
import { useGameContext } from "../../../context/GameContext"
import "./Chat.css"
import { ChatInput } from "./ChatInput/ChatInput";

enum ChatType {
    GAME_INFO = "Game info: ",
    USER_MESSAGE = "User: "
}

type ChatMessage = {
    text: string,
    type: ChatType
}

export const Chat: React.FC = () => {
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [latestUserMessage, setLatestUserMessage] = useState<string | undefined>(undefined);
    const { latestGameInfoMessage } = useGameContext();

    useEffect(() => {
        const gameInfoMessage = { text: latestGameInfoMessage, type: ChatType.GAME_INFO };
        setMessages([...messages, gameInfoMessage]);
    }, [latestGameInfoMessage])

    useEffect(() => {
        if (latestUserMessage) {
            const userMessage = { text: latestUserMessage, type: ChatType.USER_MESSAGE }; 
            setMessages([...messages, userMessage]);
        }
    }, [latestUserMessage]);

    const styledMessages = useMemo(() => {
        return messages.map((message, idx) => {
            return <div key={idx} className="chat-message">
                <span className="message-sender">{message.type}</span>
                {message.text}
            </div>
        });
    }, [messages])

    return <div className='chat'>
        <div className="chat-header underline">Chat</div>
        <div className="message-container">
            <div>
                {styledMessages}
            </div>
        </div>
        <ChatInput setLatestUserMessage={setLatestUserMessage} />
    </div>
};