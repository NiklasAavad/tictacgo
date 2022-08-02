import { useState } from "react";
import { useGameContext } from "../../../context/GameContext";
import { Board } from "./Board/Board";
import "./Game.css"
import { StartScreen } from "./Startscreen/StartScreen";

export const Game: React.FC = () => {
    const { isGameStarted } = useGameContext();

    return <div className='game'>
        {isGameStarted ? <Board/> : <StartScreen/>}
    </div>
}