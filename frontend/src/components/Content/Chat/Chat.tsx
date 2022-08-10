import React, { useCallback, useEffect, useMemo, useState } from "react";
import { useGameContext } from "../../../context/GameContext";
import { useUserContext } from "../../../context/UserContext";
import { MessageCallback, useWebsocket } from "../../../websocket/useWebsocket";
import "./Chat.css";
import { ChatInput } from "./ChatInput/ChatInput";
import { ChatMessage, InfoSender, Message } from "./Message/Message";

export const Chat: React.FC = () => {
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const { latestGameInfoMessage } = useGameContext();
    const { name } = useUserContext();

    const messageCallback: MessageCallback =  useCallback((msg: MessageEvent) => {
        const parsedMessage = JSON.parse(msg.data);
        const userMessage: ChatMessage = { sender: parsedMessage.sender, body: parsedMessage.body };
        setMessages(messages => [...messages, userMessage]);
    }, [setMessages]);

    const { connect, sendChatMessage } = useWebsocket(name, messageCallback);

    useEffect(() => {
        connect();
    }, [connect])

    useEffect(() => {
        const gameInfoMessage = { sender: InfoSender.GAME_INFO, body: latestGameInfoMessage };
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
        <ChatInput sendMessage={sendChatMessage} />
    </div>
};