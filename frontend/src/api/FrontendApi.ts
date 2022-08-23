import { Position } from "../utility/Position"

export enum GameInstruction {
    START_GAME = "start game",
    CHOOSE_SQUARE = "choose square",
    SELECT_CHARACTER = "select character",
}

// Not used yet, but should (maybe) be used to validate message before sending it
export type GameMessage = {
    instruction: GameInstruction,
    content?: Position
}