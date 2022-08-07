import { useEffect, useLayoutEffect, useRef } from "react"
import { useUserContext } from "../../context/UserContext"
import { Chat } from "./Chat/Chat"
import "./Content.css"
import { Game } from "./Game/Game"

export const Content: React.FC = () => {
    const { name, setName } = useUserContext();

    useEffect(() => {
        if (name === undefined) {
            const enteredName = prompt("Please enter your name")?.trim() || "Unknown User";
            setName(enteredName);
        }
    }, [name, setName])

    if (!name) {
        return null;
    }

    return <div className='center'>
        <Chat />
        <Game />
    </div>
}