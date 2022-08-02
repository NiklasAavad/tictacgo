import { SquareCharacter, SquareType } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";

export interface GameService {
    startGame: () => void,
    getWinningCombination: () => [Position[] | undefined, SquareCharacter],
    isGameOver: () => boolean,
    isChoiceValid: (position: Position) => boolean,
    chooseSquare: (position: Position) => SquareType | null,
    changePlayerInTurn: () => void
}