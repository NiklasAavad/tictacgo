import { useCallback } from "react"
import { Position } from "./Position"
import "./Square.css"

export type RemoveBorder = {
    top?: boolean,
    bottom?: boolean,
    left?: boolean,
    right?: boolean
}

enum RemoveBorderClass {
    TOP = "no-top-border",
    BOTTOM = "no-bottom-border",
    LEFT = "no-left-border",
    RIGHT = "no-right-border",
    SEPERATOR = " "
}

enum WinnerClass {
    WINNER = "winner",
    NO_WINNER = "no-winner"
}

export enum SquareCharacter {
    X = "X",
    O = "O",
    EMPTY = ""
}

type SquareProps = {
    isX?: boolean,
    isO?: boolean,
    removeBorder: RemoveBorder,
    position: Position,
    chooseSquare: (position: Position) => void
    winningCombination: Position[]
}

// TODO overvej om der skal lægge noget state ifht om den her square er x eller o i dette komponent i stedet for i Board.
// Board er ved at blive til et god-component
export const Square: React.FC<SquareProps> = (props) => {
    const getCharacter = useCallback((): SquareCharacter => {
        if (props.isX) return SquareCharacter.X;
        if (props.isO) return SquareCharacter.O;
        return SquareCharacter.EMPTY;
    }, [props.isX, props.isO]);

    const getBorderClasses = useCallback(() => {
        const removeBorder = props.removeBorder;
        const removeBorderClasses: RemoveBorderClass[] = []
        
        if (removeBorder.top) {
            removeBorderClasses.push(RemoveBorderClass.TOP);
        }
        if (removeBorder.bottom) {
            removeBorderClasses.push(RemoveBorderClass.BOTTOM);
        }
        if (removeBorder.left) {
            removeBorderClasses.push(RemoveBorderClass.LEFT);
        }
        if (removeBorder.right) {
            removeBorderClasses.push(RemoveBorderClass.RIGHT);
        }

        return removeBorderClasses.join(RemoveBorderClass.SEPERATOR);
    }, []);

    const getWinnerClass = useCallback((): WinnerClass => {
        const isSquareInWinningCombination = props.winningCombination.includes(props.position);
        if (isSquareInWinningCombination) {
            return WinnerClass.WINNER;
        }

        return WinnerClass.NO_WINNER;
    }, [props.winningCombination]);

    const character = getCharacter();
    const borderClasses = getBorderClasses();
    const winnerClass = getWinnerClass();

    return <div onClick={() => props.chooseSquare(props.position)} className={`square ${borderClasses} ${winnerClass}`}>
        {character}
    </div>
}