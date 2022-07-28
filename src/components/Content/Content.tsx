import { Chat } from "./Chat/Chat"
import "./Content.css"
import { Game } from "./Game/Game"

export const Content: React.FC = () => {
    return <div className='center'>
        <Chat/>
        <Game/> 
    </div>
}