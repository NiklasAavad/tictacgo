import { useState } from "react";
import { Board } from "./Board"
import "./Game.css"
import { StartScreen } from "./StartScreen";

export const Game: React.FC = () => {
    const [isGameStarted, setIsGameStarted] = useState(false);

    return <div className='game'>
        {isGameStarted ? <Board setIsSGameStarted={setIsGameStarted}/> : <StartScreen setIsGameStarted={setIsGameStarted}/>}
    </div>
}