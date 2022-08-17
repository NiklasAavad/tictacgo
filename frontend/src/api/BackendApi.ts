import { Position } from "../utility/Position";
import { SquareCharacter } from './../components/Content/Game/Square/Square';

export const BASE_URL = "localhost:8080";
export const GAME_WS_URL = `ws://${BASE_URL}/gamews`;
export const CHAT_WS_URL = `ws://${BASE_URL}/chatws`;

export enum GameCommand {
    RESULT = "result",
    GAME_OVER = "game over",
    BOARD = "board",
};

export type JSONResult = {
    WinningCombination: Position[],
    WinningCharacter: SquareCharacter,
    HasWinner: boolean
}
