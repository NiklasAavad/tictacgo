import React from "react";
import { sendMsg } from "../../../../websocket/Websocket";
import "./ChatInput.css";

export const ChatInput: React.FC = () => {
    const onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter") {
            const inputElement = event.target;
            
            const message = inputElement.value.trim();
            if (message) {
                sendMsg(message);
            }
            
            inputElement.value = "";
        }
    }

    return <div className="chat-input">
        <input placeholder="Type a message..." onKeyDown={onKeyDown} />
    </div>
}