import React, { useEffect, useMemo, useState } from "react";
import { useGameContext } from "../../../context/GameContext";
import { connect, MessageCallback } from "../../../websocket/Websocket";
import "./Chat.css";
import { ChatInput } from "./ChatInput/ChatInput";
import { ChatMessage, ChatType, Message } from "./Message/Message";

export const Chat: React.FC = () => {
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const { latestGameInfoMessage } = useGameContext();

    const messageCallback: MessageCallback = (msg: MessageEvent) => {
        const parsedMessage: ChatMessage = JSON.parse(msg.data);
        const userMessage: ChatMessage = { body: parsedMessage.body, type: ChatType.USER_MESSAGE };
        setMessages(messages => [...messages, userMessage]);
    }

    useEffect(() => {
        connect(messageCallback);
    })

    useEffect(() => {
        const gameInfoMessage = { body: latestGameInfoMessage, type: ChatType.GAME_INFO };
        setMessages(messages => [...messages, gameInfoMessage]);
    }, [latestGameInfoMessage])

    const styledMessages = useMemo(() => {
        return messages.map((message, idx) => <Message message={message} key={idx} />)
    }, [messages])

    return <div className='chat'>
        <div className="chat-header underline">Chat</div>
        <div className="message-container">
            <div>
                {styledMessages}
            </div>
        </div>
        <ChatInput />
    </div>
};