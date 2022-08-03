import "./ChatInput.css";

type ChatInputProps = {
    setLatestUserMessage: (message: string) => void
};

export const ChatInput: React.FC<ChatInputProps> = (props) => {
    const onKeyDown = (event: any) => {
        if (event.key === "Enter") {
            const message = event.target.value;
            sendMessage(message);
            event.target.value = "";
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