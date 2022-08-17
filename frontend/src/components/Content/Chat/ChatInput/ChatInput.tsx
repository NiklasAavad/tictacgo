import React from "react";
import "./ChatInput.css";

type ChatInputProps = {
    sendMessage: (msg: string) => void
}

export const ChatInput: React.FC<ChatInputProps> = (props) => {
    const onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter") {
            const inputElement = event.target;

            const message = inputElement.value.trim();
            if (message) {
                props.sendMessage(message);
            }

            inputElement.value = "";
        }
    }

    return <div className="chat-input">
        <input placeholder="Type a message..." onKeyDown={onKeyDown} />
    </div>
}