import React, { useRef } from "react";
import "./ChatInput.css";

type ChatInputProps = {
    setLatestUserMessage: (message: string) => void
};

export const ChatInput: React.FC<ChatInputProps> = (props) => {
    const onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter") {
            const inputElement = event.target;
            const message = inputElement.value;
            sendMessage(message);
            inputElement.value = "";
        }
    }

    const sendMessage = (message: string) => {
        const trimmedMessage = message.trim();
        if (trimmedMessage) {
            props.setLatestUserMessage(trimmedMessage);
        }
    }

    return <div className="chat-input">
        <input placeholder="Type a message..." onKeyDown={onKeyDown}/>
    </div>
}