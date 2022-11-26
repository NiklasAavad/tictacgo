import { useEffect, useMemo, useState } from "react";
import { useGameContext } from "../../../../context/GameContext";
import { Position } from "../../../../utility/Position";
import "./Square.css";

export type Border = {
    top?: boolean,
    bottom?: boolean,
    left?: boolean,
    right?: boolean
}

enum BorderClass {
    NO_TOP = "no-top-border",
    NO_BOTTOM = "no-bottom-border",
    NO_LEFT = "no-left-border",
    NO_RIGHT = "no-right-border",
    TOP_RIGHT_BORDER_RADIUS = "top-right-border-radius",
    BOTTOM_RIGHT_BORDER_RADIUS = "bottomright-border-radius",
}

export enum SquareCharacter {
    X = "X",
    O = "O",
}

export type EmptyString = "";

type SquareProps = {
    position: Position,
    border: Border,
}

export const Square: React.FC<SquareProps> = (props) => {
    const [character, setCharacter] = useState<SquareCharacter | EmptyString>("");
    const { board, result, chooseSquare } = useGameContext();

    useEffect(() => {
        const characterOnBoard = board[props.position];
        const squareHasBeenSelected = characterOnBoard !== character;
        if (squareHasBeenSelected) {
            setCharacter(characterOnBoard);
        }
    }, [board, props.position, character])

    const borderClasses = useMemo(() => {
        const border = props.border;
        const borderClasses: BorderClass[] = []

        if (border.top) {
            borderClasses.push(BorderClass.NO_TOP);
        }
        if (border.bottom) {
            borderClasses.push(BorderClass.NO_BOTTOM);
        }
        if (border.left) {
            borderClasses.push(BorderClass.NO_LEFT);
        }
        if (border.right) {
            borderClasses.push(BorderClass.NO_RIGHT);
        }

        if (border.top && border.right) {
            borderClasses.push(BorderClass.TOP_RIGHT_BORDER_RADIUS);
        }

        if (border.bottom && border.right) {
            borderClasses.push(BorderClass.BOTTOM_RIGHT_BORDER_RADIUS);
        }

        return borderClasses.join(" ");
    }, [props.border]);

    const winnerClass = useMemo(() => {
        const isSquareInWinningCombination = result?.winningCombination.includes(props.position);
        if (isSquareInWinningCombination) {
            return "winner";
        }

        return "";
    }, [result, props.position]);

    return <div onClick={() => chooseSquare(props.position)} className={`square ${character} ${borderClasses} ${winnerClass}`}>
        {character}
    </div>
}
