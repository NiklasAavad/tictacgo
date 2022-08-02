import { useCallback, useEffect, useMemo, useState } from "react"
import { useGameContext } from "../../../../context/GameContext"
import { Position } from "../../../../utility/Position"
import "./Square.css"

export type Border = {
    top?: boolean,
    bottom?: boolean,
    left?: boolean,
    right?: boolean
}

enum BorderClass {
    TOP = "no-top-border",
    BOTTOM = "no-bottom-border",
    LEFT = "no-left-border",
    RIGHT = "no-right-border",
    TOP_RIGHT = "top-right-border-radius",
    BOTTOM_RIGHT = "bottom-right-border-radius",
}

export enum SquareCharacter {
    X = "X",
    O = "O",
}

type ExtendedSquareCharacter = SquareCharacter | "";

export type SquareType = {
    position: Position,
    character: SquareCharacter
}

type SquareProps = {
    position: Position,
    border: Border,
}

// TODO overvej React.memo
// Problemet er, at React.memo ingen effekt har pga useGameContext, men måske man alligevel kan lave noget ala arePropsEqual og så tjek på
// latestSquare.position og winningCombination.
export const Square: React.FC<SquareProps> = (props) => {
    const [character, setCharacter] = useState<ExtendedSquareCharacter>("");
    const { latestSquare, winningCombination, chooseSquare } = useGameContext();
    
    useEffect(() => {
        const squareHasBeenSelected = latestSquare?.position === props.position;
        if (squareHasBeenSelected) {
            setCharacter(latestSquare.character);  
        }
    }, [latestSquare])

    const borderClasses = useMemo(() => {
        const border = props.border;
        const borderClasses: BorderClass[] = []
        
        if (border.top) {
            borderClasses.push(BorderClass.TOP);
        }
        if (border.bottom) {
            borderClasses.push(BorderClass.BOTTOM);
        }
        if (border.left) {
            borderClasses.push(BorderClass.LEFT);
        }
        if (border.right) {
            borderClasses.push(BorderClass.RIGHT);
        }

        if (border.top && border.right) {
            borderClasses.push(BorderClass.TOP_RIGHT);
        }

        if (border.bottom && border.right) {
            borderClasses.push(BorderClass.BOTTOM_RIGHT);
        }

        return borderClasses.join(" ");
    }, []);

    const winnerClass = useMemo(() => {
        const isSquareInWinningCombination = winningCombination?.includes(props.position);
        if (isSquareInWinningCombination) {
            return "winner";
        }

        return "";
    }, [winningCombination]);

    return <div onClick={() => chooseSquare(props.position)} className={`square ${character} ${borderClasses} ${winnerClass}`}>
        {character}
    </div>
}