import { adaptBoard, adaptWinningCharacter } from "../adapter/Adapter";
import { BASE_URL, GAME_WS } from "../api/Api";
import { Position } from "../utility/Position";
import { GameContextMutator, GameService, Result } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const socket = new WebSocket(`ws://${BASE_URL}/${GAME_WS}`);

    socket.onopen = () => {
        console.log("Connection is opened");
        sendGameMessage("Get Board");
    }

    // TODO validation!
    socket.onmessage = (msg: MessageEvent): void => {
        console.log("receiving message");
        console.log(msg);

        // TODO mangler at broadcaste at et nyt spil er gået igang
        const parsedMsg = JSON.parse(msg.data);
        const response = parsedMsg.response;
        switch (parsedMsg.command) {
            case "Result":
                console.log("Received result message");
                
                const hasWinner = response.HasWinner;
                if (!hasWinner) {
                    return;
                }

                const winningCombination = response.WinningCombination;
                const winningCharacter = adaptWinningCharacter(response.WinningCharacter);
                const result = {winningCombination, winningCharacter};
                
                return gameContextMutator.setResult(result);
            case "Game Over":
                console.log("Received game over message");
                return gameContextMutator.setIsGameOver(true);
            case "Board":
                gameContextMutator.setIsGameStarted(true);
                const backendBoard = response;
                const adaptedBoard = adaptBoard(backendBoard);
                return gameContextMutator.setBoard(adaptedBoard);
            case "Player In Turn":
                console.log("Received Player In Turn message");
                return;
        }
    };

    const sendGameMessage = (instruction: string, content?: Position) => {
        console.log("Sending game message");

        const jsonMessage = JSON.stringify({ instruction, content })

        socket.send(jsonMessage);
    }

    const startGame = (): void => {
        console.log("Starting game")
        sendGameMessage("Start Game");
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