import { JSONResult, ResponseType } from "../api/BackendApi";
import { GameInstruction } from "../api/FrontendApi";
import { GameInfoMessage } from "../context/GameContext";
import { Position } from "../utility/Position";
import { GAME_WS_URL, JSONWeclome } from './../api/BackendApi';
import { SquareCharacter } from './../components/Content/Game/Square/Square';
import { getEmptyBoard } from './../utility/GameServiceUtility';
import { Board, GameContextMutator, GameService } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const socket = new WebSocket(GAME_WS_URL); // TOOD husk at brug username.

    // TODO validation!
    socket.onmessage = (msg: MessageEvent): void => {
        console.log(msg);
        const { responseType, body } = JSON.parse(msg.data);
        switch (responseType) {
            case ResponseType.RESULT:
                return resultDidChange(body)
            case ResponseType.GAME_OVER:
                return gameDidEnd();
            case ResponseType.BOARD:
                return boardDidChange(body);
            case ResponseType.NEW_MESSAGE:
                return newGameMessageReceived(body);
            case ResponseType.CHARACTER_SELECTED:
                return characterHasBeenSelected(body);
            case ResponseType.GAME_STARTED:
                return gameDidStart();
            case ResponseType.WELCOME:
                return welcomeMessageWasReceived(body);
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
        gameContextMutator.setBoard(board);
    }

    const newGameMessageReceived = (message: string) => {
        console.log("New game message received", message);
    }

    const characterHasBeenSelected = (character: SquareCharacter) => {
        if (character === SquareCharacter.X) {
            gameContextMutator.setLatestGameInfoMessage(GameInfoMessage.X_SELECTED);
        } else {
            gameContextMutator.setLatestGameInfoMessage(GameInfoMessage.O_SELECTED);
        }
    }

    const gameDidStart = () => {
        gameContextMutator.setBoard(getEmptyBoard());
        gameContextMutator.setLatestGameInfoMessage(GameInfoMessage.NEW_GAME_STARTED);
        gameContextMutator.setIsGameStarted(true);
        gameContextMutator.setResult(undefined);
        gameContextMutator.setIsGameOver(false);
    }

    // TODO fix any
    const welcomeMessageWasReceived = ({ isGameStarted, board, xClient, oClient }: JSONWeclome) => {
        gameContextMutator.setIsGameStarted(isGameStarted)

        if (isGameStarted) {
            gameContextMutator.setBoard(board);
        }

        if (xClient) {
            console.log(xClient)
        }

        if (oClient) {
            console.log(oClient)
        }
    }

    const sendGameMessage = (instruction: GameInstruction, content?: Position | SquareCharacter) => {
        const jsonMessage = JSON.stringify({ instruction, content })
        socket.send(jsonMessage);
    }

    const startGame = (): void => {
        sendGameMessage(GameInstruction.START_GAME);
    };

    const chooseSquare = (position: Position): void => {
        sendGameMessage(GameInstruction.CHOOSE_SQUARE, position);
    };

    const selectCharacter = (character: SquareCharacter): void => {
        sendGameMessage(GameInstruction.SELECT_CHARACTER, character);
    }

    return {
        startGame,
        chooseSquare,
        selectCharacter
    }
}

export default OnlineMultiplayerGameService;