import "./StartScreen.css"

type StartScreenprops = {
    setIsGameStarted: (flag: boolean) => void
}

export const StartScreen: React.FC<StartScreenprops> = (props) => {
    return <div className="start-screen-container">
        <div className="start-screen-header">Start a new game!</div>
        <button onClick={() => props.setIsGameStarted(true)} className="new-game-button">New Game</button>
    </div>
}