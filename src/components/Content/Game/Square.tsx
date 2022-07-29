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
}

// TODO overvej om der skal lægge noget state ifht om den her square er x eller o i dette komponent i stedet for i Board.
// Board er ved at blive til et god-component
export const Square: React.FC<SquareProps> = (props) => {
    const getCharacter = useCallback(() => {
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

    const character = getCharacter();
    const borderClasses = getBorderClasses();

    return <div onClick={() => props.chooseSquare(props.position)} className={"square" + RemoveBorderClass.SEPERATOR + borderClasses}>
        {character}
    </div>
}