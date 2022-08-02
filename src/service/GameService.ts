import { SquareCharacter, SquareType } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";

export type Result = {
    winningCombination: Position[],
    winningCharacter: SquareCharacter;
}

export interface GameService {
    startGame: () => void,
    getResult: () => Result | undefined,
    isGameOver: () => boolean,
    isChoiceValid: (position: Position) => boolean,
    chooseSquare: (position: Position) => SquareType,
    changePlayerInTurn: () => void
}