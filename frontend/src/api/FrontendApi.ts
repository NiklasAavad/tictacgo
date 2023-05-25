import { Position } from "../utility/Position"

export const GameInstruction = {
    START_GAME: "start game",
    CHOOSE_SQUARE: "choose square",
    SELECT_CHARACTER: "select character",
} as const;

export type GameInstructionType = typeof GameInstruction[keyof typeof GameInstruction];

// Not used yet, but should (maybe) be used to validate message before sending it
export type GameMessage = {
    instruction: GameInstructionType,
    content?: Position
}
