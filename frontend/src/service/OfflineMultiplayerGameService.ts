import { SquareCharacter } from "../components/Content/Game/Square/Square";
import { getEmptyBoard, hasPlayerWon } from "../utility/GameServiceUtility";
import { allPositions, Position } from "../utility/Position";
import { WINNING_COMBINATIONS } from "../utility/WinningCombinations";
import { GameContextMutator, GameService, Result } from "./GameService";

const OfflineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    let x: Position[] = []; // TODO måske udfases?
    let o: Position[] = []; // TODO måske udfases?
    let board = getEmptyBoard();
    let playerInTurn: SquareCharacter = SquareCharacter.X;

    const startGame = (): void => {
        x = [];
        o = [];
        board = getEmptyBoard();
        playerInTurn = SquareCharacter.X;
    };

    const getResult = (): Result | undefined => {
        const xWinningCombination = WINNING_COMBINATIONS.find(combination => {
            return combination.every(position => x.includes(position));
        });

        if (xWinningCombination) {
            return { winningCombination: xWinningCombination, winningCharacter: SquareCharacter.X }
        }

        const oWinningCombination = WINNING_COMBINATIONS.find(combination => {
            return combination.every(position => o.includes(position));
        });

        if (oWinningCombination) {
            return { winningCombination: oWinningCombination, winningCharacter: SquareCharacter.O }
        }

        return undefined;
    };

    const isGameOver = (): boolean => {
        const notEnoughInputs = x.length < 3;
        if (notEnoughInputs) {
            return false;
        }

        const allPositionsOccupied = x.length + o.length === allPositions.length;
        if (allPositionsOccupied) {
            return true;
        }

        return hasPlayerWon(x) || hasPlayerWon(o);
    };

    const isChoiceValid = (position: Position): boolean => {
        if (isGameOver()) {
            return false;
        }

        const isAlreadyX = x.includes(position);
        const isAlreadyO = o.includes(position);
        const isPositionOccupied = isAlreadyX || isAlreadyO;

        return !isPositionOccupied;
    };

    const addSquareToBoard = (position: Position): void => {
        if (playerInTurn === SquareCharacter.X) {
            x.push(position);
        } else {
            o.push(position);
        }
        board[position] = playerInTurn;
        gameContextMutator.setBoard([...board]);
    }

    const chooseSquare = (position: Position): void => {
        if (!isChoiceValid(position)) {
            return;
        }

        addSquareToBoard(position);

        if (isGameOver()) {
            const result = getResult();
            gameContextMutator.setResult(result);
            gameContextMutator.setIsGameOver(true)
        };

        changePlayerInTurn();
    };

    const changePlayerInTurn = (): void => {
        if (playerInTurn === SquareCharacter.X) {
            playerInTurn = SquareCharacter.O;
        } else {
            playerInTurn = SquareCharacter.X;
        }
    };

    return {
        startGame,
        getResult,
        isGameOver,
        isChoiceValid,
        chooseSquare,
        changePlayerInTurn
    }
}

export default OfflineMultiplayerGameService;