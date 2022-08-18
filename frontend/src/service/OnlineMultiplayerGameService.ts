import { GameCommand, JSONResult } from "../api/BackendApi";
import { GameInstruction } from "../api/FrontendApi";
import { Position } from "../utility/Position";
import { GAME_WS_URL } from './../api/BackendApi';
import { Board, GameContextMutator, GameService } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const socket = new WebSocket(GAME_WS_URL); // TOOD husk at brug username.

    socket.onopen = () => {
        sendGameMessage(GameInstruction.GET_BOARD);
    }

    // TODO validation!
    socket.onmessage = (msg: MessageEvent): void => {
        console.log(msg);
        const { command, body } = JSON.parse(msg.data);
        switch (command) {
            case GameCommand.RESULT:
                return resultDidChange(body)
            case GameCommand.GAME_OVER:
                return gameDidEnd();
            case GameCommand.BOARD:
                return boardDidChange(body);
            case GameCommand.NEW_MESSAGE:
                return newGameMessageReceived(body);

        }
        throw new Error("No command matched the received message: " + JSON.stringify(msg));
    };


    const resultDidChange = (jsonResult: JSONResult) => {
        if (!jsonResult.HasWinner) {
            return;
        }

        const winningCombination = jsonResult.WinningCombination;
        const winningCharacter = jsonResult.WinningCharacter;
        const result = { winningCombination, winningCharacter };

        gameContextMutator.setResult(result);
    }

    const gameDidEnd = () => {
        gameContextMutator.setIsGameOver(true);
    }

    const boardDidChange = (board: Board) => {
        // TODO overvej om de tre nedenstående metoder skal trækkes over i en "gameDidStart" command. 
        gameContextMutator.setIsGameStarted(true);
        gameContextMutator.setIsGameOver(false);
        gameContextMutator.setResult(undefined);

        gameContextMutator.setBoard(board);
    }

    const newGameMessageReceived = (message: string) => {
        console.log("New game message received", message);
    }

    const sendGameMessage = (instruction: GameInstruction, content?: Position) => {
        const jsonMessage = JSON.stringify({ instruction, content })
        socket.send(jsonMessage);
    }

    const startGame = (): void => {
        sendGameMessage(GameInstruction.START_GAME);
    };

    const chooseSquare = (position: Position): void => {
        sendGameMessage(GameInstruction.CHOOSE_SQUARE, position);
    };

    return {
        startGame,
        chooseSquare,
    }
}

export default OnlineMultiplayerGameService;