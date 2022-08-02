import { GameProvider } from "../../context/GameContext"
import { GameService } from "../../service/GameService"
import OfflineMultiplayerGameService from "../../service/OfflineMultiplayerGameService"
import { Chat } from "./Chat/Chat"
import "./Content.css"
import { Game } from "./Game/Game"

export const Content: React.FC = () => {
    const gameService: GameService = OfflineMultiplayerGameService;

    return <div className='center'>
        <Chat/>
        <GameProvider gameService={gameService}>
            <Game/>
        </GameProvider>
    </div>
}