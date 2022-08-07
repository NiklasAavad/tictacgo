import "./Message.css";

export enum ChatType {
    GAME_INFO = "Game info: ",
    USER_MESSAGE = "User: "
}

export type ChatMessage = {
    type: ChatType
    body: string,
}

type MessageProps = {
    message: ChatMessage,
    key: number
}

export const Message: React.FC<MessageProps> = (props) => {
    return <div key={props.key} className="chat-message">
        <span className="message-sender">{props.message.type}</span>
        {props.message.body}
    </div>
}