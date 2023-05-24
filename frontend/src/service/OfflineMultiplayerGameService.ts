import { SquareCharacter } from "../components/Content/Game/Square/Square";
import { getEmptyBoard, startGameForMutator } from "../utility/GameServiceUtility";
import { allPositions, Position } from "../utility/Position";
import { WINNING_COMBINATIONS } from "../utility/WinningCombinations";
import { Board, GameContextMutator, GameService, Result } from "./GameService";

const OfflineMultiplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    let board: Board = getEmptyBoard();
    let playerInTurn: SquareCharacter = SquareCharacter.X;

    const startGame = (): void => {
        board = getEmptyBoard();
        playerInTurn = SquareCharacter.X;
		startGameForMutator(gameContextMutator);
    };

    const getResult = (): Result | undefined => {
        for (const winningCombination of WINNING_COMBINATIONS) {
            if (isWinningRow(winningCombination)) {
                const winningCharacter = board[winningCombination[0]];

                if (winningCharacter === "") {
                    throw new Error("Game had a result, but no winning character");
                }

                return { winningCombination, winningCharacter };
            }
        }

        return undefined;
    };

    const isWinningRow = (row: Position[]) => {
        const isNotEmptyRow = board[row[0]] !== ""
        const isSameCharacter = board[row[0]] === board[row[1]] && board[row[1]] === board[row[2]];
        return isNotEmptyRow && isSameCharacter;
    }

    const hasWinner = (): boolean => {
        return WINNING_COMBINATIONS.some(combination => isWinningRow(combination));
    }

    const getOccupiedSquares = () => {
        return board.filter(square => square !== "");
    }

    const isGameOver = (): boolean => {
        const occupiedSquares = getOccupiedSquares();
        const notEnoughInputs = occupiedSquares.length < 5;
        if (notEnoughInputs) {
            return false;
        }

        const allPositionsOccupied = occupiedSquares.length === allPositions.length;
        if (allPositionsOccupied) {
            return true;
        }

        return hasWinner();
    };

    const isChoiceValid = (position: Position): boolean => {
        if (isGameOver()) {
            return false;
        }

        const isSquareFree = board[position] === "";
        return isSquareFree;
    };

    const addSquareToBoard = (position: Position): void => {
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

    const selectCharacter = (character: SquareCharacter): void => {
        console.log("Character was selected", character);
    }

    return {
        startGame,
        getResult,
        isGameOver,
        isChoiceValid,
        chooseSquare,
        changePlayerInTurn,
        selectCharacter
    }
}

export default OfflineMultiplayerGameService;
