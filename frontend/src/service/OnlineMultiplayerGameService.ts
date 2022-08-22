import { JSONResult, ResponseType } from "../api/BackendApi";
import { GameInstruction } from "../api/FrontendApi";
import { GameInfoMessage } from "../context/GameContext";
import { Position } from "../utility/Position";
import { GAME_WS_URL } from './../api/BackendApi';
import { SquareCharacter } from './../components/Content/Game/Square/Square';
import { getEmptyBoard } from './../utility/GameServiceUtility';
import { Board, GameContextMutator, GameService } from "./GameService";

const OnlineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const socket = new WebSocket(GAME_WS_URL); // TOOD husk at brug username.

    socket.onopen = () => {
        sendGameMessage(GameInstruction.GET_BOARD);
    }

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