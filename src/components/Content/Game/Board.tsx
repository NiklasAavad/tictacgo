import { useState } from "react"
import "./Board.css"
import { Position } from "./Position"
import { Border, Square, SquareCharacter } from "./Square"

export const Board: React.FC = () => {
    const [x, setX] = useState<Position[]>([]);
    const [o, setO] = useState<Position[]>([]);
    const [playerInTurn, setPlayerInTurn] = useState<SquareCharacter>(SquareCharacter.X);

    const isChoiceValid = (position: Position) => {
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
    const createSquare = (position: Position, border: Border = {}) => {
        const isX = x.includes(position);
        const isO = o.includes(position);

        return <Square isX={isX} isO={isO} border={border} position={position} chooseSquare={chooseSquare}/>
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