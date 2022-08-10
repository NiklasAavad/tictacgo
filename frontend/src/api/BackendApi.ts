import { BackendSquareChacter, BackendWinningCharacter } from "../adapter/Adapter";
import { Position } from "../utility/Position";

export const BASE_URL = "localhost:8080";
export const GAME_WS = "gamews";
export const CHAT_WS = "chatws";

export enum GameCommand {
    RESULT = "Result",
    GAME_OVER = "Game Over",
    BOARD = "Board",
    PLAYER_IN_TURN = "Player In Turn"
};

export type JSONResult = {
    WinningCombination: Position[],
    WinningCharacter: BackendWinningCharacter,
    HasWinner: boolean
}

export type GameBody = JSONResult | BackendSquareChacter[] | BackendSquareChacter | boolean

// Not used yet, but should be used to validate the response from the game.
export type GameResponse = {
    command: GameCommand,
    response: GameBody
};
