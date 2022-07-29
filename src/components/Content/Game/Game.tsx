import { useState } from "react";
import { Board } from "./Board"
import "./Game.css"
import { StartScreen } from "./StartScreen";

export const Game: React.FC = () => {
    const [isStarted, setIsStarted] = useState(true);


    return <div className='game'>
        {isStarted ? <Board/> : <StartScreen/>}
    </div>
}