import "./Message.css";

export enum InfoSender {
    GAME_INFO = "Game Info",
    CHAT_INFO = "Chat Info"
}

export type ChatMessage = {
    sender: string,
    body: string,
}

type MessageProps = {
    message: ChatMessage,
    key: number
}

export const Message: React.FC<MessageProps> = (props) => {
    return <div key={props.key} className="chat-message">
        <span className="message-sender">{props.message.sender}: </span>
        {props.message.body}
    </div>
}