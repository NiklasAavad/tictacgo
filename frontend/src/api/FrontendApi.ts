import { Position } from "../utility/Position"

export enum GameInstruction {
    START_GAME = "start game",
    CHOOSE_SQUARE = "choose square",
    GET_BOARD = "get board"
}

// Not used yet, but should (maybe) be used to validate message before sending it
export type GameMessage = {
    instruction: GameInstruction,
    content?: Position
}