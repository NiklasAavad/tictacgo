import { JSONResult, ResponseType } from "../api/BackendApi";
import { GameInstruction, GameInstructionType, GameMessage } from "../api/FrontendApi";
import { GAMEINFO } from "../context/GameContext";
import { Position } from "../utility/Position";
import { GAME_WS_URL, JSONWeclome } from './../api/BackendApi';
import { SquareCharacter } from './../components/Content/Game/Square/Square';
import { startGameForMutator } from './../utility/GameServiceUtility';
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
            gameContextMutator.setLatestGameInfoMessage(GAMEINFO.X_SELECTED);
        } else {
            gameContextMutator.setLatestGameInfoMessage(GAMEINFO.O_SELECTED);
        }
    }

    const gameDidStart = () => {
		startGameForMutator(gameContextMutator);
    }

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

    const sendGameMessage = (msg: GameMessage) => {
        const jsonMessage = JSON.stringify(msg)
        socket.send(jsonMessage);
    }

    const startGame = (): void => {
		const msg = { instruction: GameInstruction.START_GAME };
        sendGameMessage(msg);
    };

    const chooseSquare = (position: Position): void => {
		const msg = {
			instruction: GameInstruction.CHOOSE_SQUARE,
			content: position
		}
        sendGameMessage(msg);
    };

    const selectCharacter = (character: SquareCharacter): void => {
		const msg = {
			instruction: GameInstruction.SELECT_CHARACTER,
			content: character
		}
        sendGameMessage(msg);
    }

    return {
        startGame,
        chooseSquare,
        selectCharacter
    }
}

export default OnlineMultiplayerGameService;
