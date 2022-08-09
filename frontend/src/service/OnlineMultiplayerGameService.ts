import { adaptBoard, adaptWinningCharacter, BackendSquareChacter, BackendWinningCharacter } from "../adapter/Adapter";
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
    // TODO Broadcast et nyt spil er gået igang!
    socket.onmessage = (msg: MessageEvent): void => {
        const parsedMsg = JSON.parse(msg.data);
        const response = parsedMsg.response;
        switch (parsedMsg.command) {
            case "Result":
                return resultDidChange(response)
            case "Game Over":
                return gameDidEnd();
            case "Board":
                return boardDidChange(response);
            case "Player In Turn":
                return playerInTurnDidChange(response);
        }
        throw new Error("No command matched the received message: " + msg);
    };

    type JSONResult = {
        WinningCombination?: Position[],
        WinningCharacter?: BackendWinningCharacter,
        HasWinner: boolean
    }

    const resultDidChange = (jsonResult: JSONResult) => {
        if (!jsonResult.HasWinner) {
            return;
        }

        console.log(jsonResult.WinningCharacter);

        const winningCombination = jsonResult.WinningCombination!;
        const winningCharacter = adaptWinningCharacter(jsonResult.WinningCharacter!);
        const result = { winningCombination, winningCharacter };

        gameContextMutator.setResult(result);
    }

    const gameDidEnd = () => {
        gameContextMutator.setIsGameOver(true);
    }

    const boardDidChange = (board: BackendSquareChacter[]) => {
        gameContextMutator.setIsGameStarted(true);
        const adaptedBoard = adaptBoard(board);
        gameContextMutator.setBoard(adaptedBoard);
    }

    // Unused for now
    const playerInTurnDidChange = (playerInTurn: BackendSquareChacter) => {
        // Do nothing
    }

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