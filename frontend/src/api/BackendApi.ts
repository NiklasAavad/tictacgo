import { Board } from "../service/GameService";
import { Position } from "../utility/Position";
import { SquareCharacter } from './../components/Content/Game/Square/Square';

export const BASE_URL = "localhost:8080";
export const GAME_WS_URL = `ws://${BASE_URL}/gamews`;
export const CHAT_WS_URL = `ws://${BASE_URL}/chatws`;

export enum ResponseType {
    RESULT = "result",
    GAME_OVER = "game over",
    BOARD = "board",
    NEW_MESSAGE = "new message",
    CHARACTER_SELECTED = "character selected",
    GAME_STARTED = "game started",
    WELCOME = "welcome"
};

export type JSONResult = {
    WinningCombination: Position[],
    WinningCharacter: SquareCharacter,
    HasWinner: boolean
}

export type JSONWeclome = {
    isGameStarted: boolean,
    board: Board,
    xClient: string,
    oClient: string
}
