import { Position } from "../utility/Position"

export enum GameInstruction {
    START_GAME = "Start Game",
    CHOOSE_SQUARE = "Choose Square",
    GET_BOARD = "Get Board"
}

// Not used yet, but should (maybe) be used to validate message before sending it
export type GameMessage = {
    instruction: GameInstruction,
    content?: Position
}