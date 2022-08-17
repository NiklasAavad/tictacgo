import React from "react";
import { Position } from "../../../../utility/Position";
import { Border, Square } from "../Square/Square";
import "./Board.css";

export const Board: React.FC = React.memo(() => {
    const createSquare = (position: Position, border: Border) => {
        return <Square key={position} position={position} border={border} />
    }

    const createTopSquares = () => {
        const topLeftBorders = { top: true, left: true };
        const topCenterBorders = { top: true };
        const topRightBorders = { top: true, right: true };

        const topLeftSquare = createSquare(Position.TOP_LEFT, topLeftBorders);
        const topCenterSquare = createSquare(Position.TOP_CENTER, topCenterBorders)
        const topRightSquare = createSquare(Position.TOP_RIGHT, topRightBorders);

        return [topLeftSquare, topCenterSquare, topRightSquare];
    };

    const createCenterSquares = () => {
        const centerLeftBorders = { left: true };
        const centerBorders = {};
        const centerRightBorders = { right: true };

        const centerLeftSquare = createSquare(Position.CENTER_LEFT, centerLeftBorders);
        const centerSquare = createSquare(Position.CENTER, centerBorders)
        const centerRightSquare = createSquare(Position.CENTER_RIGHT, centerRightBorders);

        return [centerLeftSquare, centerSquare, centerRightSquare]
    }

    const createBottomSquares = () => {
        const bottomLeftBorders = { bottom: true, left: true };
        const bottomCenterBorders = { bottom: true };
        const bottomRightBorders = { bottom: true, right: true };

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
});