import { useEffect, useState } from "react"
import { Position } from "../../../../utility/Position"
import { WINNING_COMBINATIONS } from "../../../../utility/WinningCombinations"
import { RemoveBorder, Square, SquareCharacter } from "../Square/Square"
import "./Board.css"

type BoardProps = {
    setIsSGameStarted: (flag: boolean) => void
}

export const Board: React.FC<BoardProps> = (props) => {
    const [x, setX] = useState<Position[]>([]);
    const [o, setO] = useState<Position[]>([]);
    const [playerInTurn, setPlayerInTurn] = useState<SquareCharacter>(SquareCharacter.X);
    const [winningCombination, setWinningCombination] = useState<Position[] | null>(null)

    useEffect(() => {
        if (isGameOver()) {
            endGame();
        }
    }, [x, o])

    const endGame = () => {
        const [winningCombination, winningCharacter] = getWinningCombination();

        if (winningCombination) {
            console.log(`Game was won by ${winningCharacter}!`)
            setWinningCombination(winningCombination)
        } else {
            console.log("Game was tied...");
        }

        const waitForNewGame = setTimeout(() => props.setIsSGameStarted(false), 2500)
        return () => clearTimeout(waitForNewGame);
    }

    const getWinningCombination = (): [Position[] | undefined, SquareCharacter] => {
        const xWinningCombination = WINNING_COMBINATIONS.find(combination => {
            return combination.every(position => x.includes(position))
        });
        
        if (xWinningCombination) {
            return [xWinningCombination, SquareCharacter.X];
        }

        const oWinningCombination = WINNING_COMBINATIONS.find(combination => {
            return combination.every(position => o.includes(position))
        });

        if (oWinningCombination) {
            return [oWinningCombination, SquareCharacter.O]
        }

        return [undefined, SquareCharacter.EMPTY];
    }

    const hasPlayerWon = (playerPositions: Position[]) => {
        return WINNING_COMBINATIONS.some(combination => {
            return combination.every(position => playerPositions.includes(position))
        });
    }

    const isGameOver = () => {
        const notEnoughInputs = x.length < 3;
        if (notEnoughInputs) {
            return false;
        }

        const allPositionsOccupied = x.length + o.length == 9; // TODO magisk konstant? 
        if (allPositionsOccupied) {
            return true;
        }

        return hasPlayerWon(x) || hasPlayerWon(o); 
    }

    const isChoiceValid = (position: Position) => {
        if (isGameOver()) {
            return false;
        }

        const isAlreadyX = x.includes(position);
        const isAlreadyO = o.includes(position)
        const isPositionOccupied = isAlreadyX || isAlreadyO

        return !isPositionOccupied;
    }

    const chooseSquare = (position: Position) => {
        if (!isChoiceValid(position)) {
            return;
        }

        if (playerInTurn === SquareCharacter.X) {
            const newX = [...x, position];
            setX(newX);
            setPlayerInTurn(SquareCharacter.O);
        } else {
            const newO = [...o, position];
            setO(newO);
            setPlayerInTurn(SquareCharacter.X);
        }
    }
    
    // TODO overvej om det her kan trækkes ud i en en service eller andet. Her skal createSquare dog bruge x og o, hvilket bliver til mange argumenter.
    const createSquare = (position: Position, removeBorder: RemoveBorder) => {
        const isX = x.includes(position);
        const isO = o.includes(position);

        return <Square key={position} isX={isX} isO={isO} removeBorder={removeBorder} position={position} chooseSquare={chooseSquare} winningCombination={winningCombination}/>
    }

    const createTopSquares = () => {
        const topLeftBorders = {top: true, left: true};
        const topCenterBorders = {top: true};
        const topRightBorders = {top: true, right: true};

        const topLeftSquare = createSquare(Position.TOP_LEFT, topLeftBorders);
        const topCenterSquare = createSquare(Position.TOP_CENTER, topCenterBorders)
        const topRightSquare = createSquare(Position.TOP_RIGHT, topRightBorders);

        return [topLeftSquare, topCenterSquare, topRightSquare];
    }
    
    const createCenterSquares = () => {
        const centerLeftBorders = {left: true};
        const centerBorders = {};
        const centerRightBorders = {right: true};

        const centerLeftSquare = createSquare(Position.CENTER_LEFT, centerLeftBorders);
        const centerSquare = createSquare(Position.CENTER, centerBorders)
        const centerRightSquare = createSquare(Position.CENTER_RIGHT, centerRightBorders);

        return [centerLeftSquare, centerSquare, centerRightSquare]
    }

    const createBottomSquares = () => {
        const bottomLeftBorders = {bottom: true, left: true};
        const bottomCenterBorders = {bottom: true};
        const bottomRightBorders = {bottom: true, right: true};

        const bottomLeftSquare = createSquare(Position.BOTTOM_LEFT, bottomLeftBorders);
        const bottomCenterSquare = createSquare(Position.BOTTOM_CENTER, bottomCenterBorders);
        const bottomRightSquare = createSquare(Position.BOTTOM_RIGHT, bottomRightBorders);

        return [bottomLeftSquare, bottomCenterSquare, bottomRightSquare];
    }

    const topSquares = createTopSquares();
    const centerSquares = createCenterSquares();
    const bottomSquares = createBottomSquares();

    return <div className="board">
        <div className="row">
            {topSquares}
        </div>
        <div className="row">
            {centerSquares}
        </div>
        <div className="row">
            {bottomSquares}
        </div>
    </div>
}