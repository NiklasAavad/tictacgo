import { BASE_URL, GAME_WS } from "../api/Api";
import { Position } from "../utility/Position";
import { GameContextMutator, GameService, Result } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const socket = new WebSocket(`ws://${BASE_URL}/${GAME_WS}`);

    socket.onmessage = (msg: MessageEvent) => {
        console.log("receiving message");
        console.log(msg);
    };

    const sendGameMessage = (instruction: string, content?: Position) => {
        console.log("Sending game message");

        const jsonMessage = JSON.stringify({ instruction, content })

        socket.send(jsonMessage);
    }

    const startGame = (): void => {
        console.log("Starting game")
    };

    const getResult = (): Result | undefined => {
        console.log("Getting result")
        return undefined;
    };

    const isGameOver = (): boolean => {
        console.log("Checking if game is over")
        return false;
    };

    const isChoiceValid = (position: Position): boolean => {
        console.log("Checking if choice is valid")
        return true;
    };

    const chooseSquare = (position: Position): void => {
        console.log("Choosing square");
        sendGameMessage("Choose Square", position); // TODO adapter for position! Enten her eller på backend.
    };

    const changePlayerInTurn = (): void => {
        console.log("Changing player in turn")
        sendGameMessage("Change Player In Turn")
    };

    return {
        startGame,
        getResult,
        isGameOver,
        isChoiceValid,
        chooseSquare,
        changePlayerInTurn
    }
}

export default OnlineMultiplayerGameService;