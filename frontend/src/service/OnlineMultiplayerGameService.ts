import { adaptBoard, adaptWinningCharacter, BackendSquareChacter } from "../adapter/Adapter";
import { BASE_URL, GameCommand, GAME_WS, JSONResult } from "../api/BackendApi";
import { GameInstruction } from "../api/FrontendApi";
import { Position } from "../utility/Position";
import { GameContextMutator, GameService } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const socket = new WebSocket(`ws://${BASE_URL}/${GAME_WS}`);

    socket.onopen = () => {
        // sendGameMessage(GameInstruction.GET_BOARD);
    }

    // TODO validation!
    socket.onmessage = (msg: MessageEvent): void => {
        const { command, body } = JSON.parse(msg.data);
        switch (command) {
            case GameCommand.RESULT:
                return resultDidChange(body)
            case GameCommand.GAME_OVER:
                return gameDidEnd();
            case GameCommand.BOARD:
                return boardDidChange(body);
            case GameCommand.PLAYER_IN_TURN:
                return playerInTurnDidChange(body);
        }
        throw new Error("No command matched the received message: " + JSON.stringify(msg));
    };


    const resultDidChange = (jsonResult: JSONResult) => {
        if (!jsonResult.HasWinner) {
            return;
        }

        const winningCombination = jsonResult.WinningCombination;
        const winningCharacter = adaptWinningCharacter(jsonResult.WinningCharacter);
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