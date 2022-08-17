import { useGameContext } from "../../../../context/GameContext";
import "./StartScreen.css";

export const StartScreen: React.FC = () => {
    const { startGame } = useGameContext();

    return <div className="start-screen-container">
        <div className="start-screen-header">Start a new game!</div>
        <button onClick={startGame} className="new-game-button">New Game</button>
    </div>
}