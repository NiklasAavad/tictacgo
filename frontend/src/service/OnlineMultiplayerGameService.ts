import { Position } from "../utility/Position";
import { MessageCallback, useWebsocket } from "../websocket/useWebsocket";
import { GameContextMutator, GameService, Result } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const { connect, sendGameMessage } = useWebsocket("Test", true);

    const messageCallback: MessageCallback = (msg: MessageEvent) => {
        const parsedMessage = JSON.parse(msg.data);
        const board = parsedMessage.board;
        console.log(board);
    }

    connect(messageCallback);

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