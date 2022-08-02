import { SquareCharacter, SquareType } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";
import { WINNING_COMBINATIONS } from "../utility/WinningCombinations";
import { GameService } from "./GameService";

let x: Position[] = [];
let o: Position[] = [];
let playerInTurn: SquareCharacter = SquareCharacter.X;
const MAX_POSITIONS = 9;

const OfflineMultiplayerGameService: GameService = {
    startGame: function (): void {
        x = [];
        o = [];
        playerInTurn = SquareCharacter.X;
    },

    getWinningCombination: function (): [Position[] | undefined, SquareCharacter] {
        const xWinningCombination = WINNING_COMBINATIONS.find(combination => {
            return combination.every(position => x.includes(position));
        });

        if (xWinningCombination) {
            return [xWinningCombination, SquareCharacter.X];
        }

        const oWinningCombination = WINNING_COMBINATIONS.find(combination => {
            return combination.every(position => o.includes(position));
        });

        if (oWinningCombination) {
            return [oWinningCombination, SquareCharacter.O];
        }

        return [undefined, SquareCharacter.EMPTY];
    },

    isGameOver: function (): boolean {
        const notEnoughInputs = x.length < 3;
        if (notEnoughInputs) {
            return false;
        }

        const allPositionsOccupied = x.length + o.length == MAX_POSITIONS;
        if (allPositionsOccupied) {
            return true;
        }

        return hasPlayerWon(x) || hasPlayerWon(o);
    },

    isChoiceValid: function (position: Position): boolean {
        if (OfflineMultiplayerGameService.isGameOver()) {
            return false;
        }

        const isAlreadyX = x.includes(position);
        const isAlreadyO = o.includes(position);
        const isPositionOccupied = isAlreadyX || isAlreadyO;

        return !isPositionOccupied;
    },

    chooseSquare: function (position: Position): SquareType {
        if (playerInTurn === SquareCharacter.X) {
            x.push(position);
        } else {
            o.push(position);
        }

        return { position: position, character: playerInTurn };
    },

    changePlayerInTurn: function (): void {
        if (playerInTurn == SquareCharacter.X) {
            playerInTurn = SquareCharacter.O;
        } else {
            playerInTurn = SquareCharacter.X;
        }
    },
}

const hasPlayerWon = (playerPositions: Position[]) => {
    return WINNING_COMBINATIONS.some(combination => {
        return combination.every(position => playerPositions.includes(position))
    });
}

export default OfflineMultiplayerGameService;