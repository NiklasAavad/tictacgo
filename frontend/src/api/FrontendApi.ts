import { SquareCharacter } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position"

export const GameInstruction = {
    START_GAME: "start game",
    CHOOSE_SQUARE: "choose square",
    SELECT_CHARACTER: "select character",
} as const;

export type GameInstructionType = typeof GameInstruction[keyof typeof GameInstruction];

// TODO consider creating discriminated union type for GameMessage. 
// This would allow to validate message before sending it
// https://www.typescriptlang.org/docs/handbook/advanced-types.html#discriminated-unions
// However, might be hard to implement properly on the backend?
export type GameMessage = {
    instruction: GameInstructionType,
    content?: Position | SquareCharacter
}
