import { useGameContext } from "../../../../context/GameContext";
import { SquareCharacter } from "../Square/Square";
import "./StartScreen.css";

export const StartScreen: React.FC = () => {
    const { startGame, selectCharacter } = useGameContext();

    return <div className="start-screen-container">
        <div className="start-screen-header">Start a new game!</div>
        <button onClick={() => selectCharacter(SquareCharacter.X)} className="select-character">
            Choose X
        </button>
        <button onClick={() => selectCharacter(SquareCharacter.O)} className="select-character">
            Choose O
        </button>
        <div className="divider" />
        <button onClick={startGame} className="new-game-button">
            Start Game
        </button>
    </div>
}