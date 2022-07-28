import { Position } from "./Position"
import "./Square.css"

export type Border = {
    top?: boolean,
    bottom?: boolean,
    left?: boolean,
    right?: boolean
}

type SquareProps = {
    isX?: boolean,
    isO?: boolean,
    border: Border,
    position: Position,
    chooseSquare: (position: Position) => void
}

export enum SquareCharacter {
    X = "X",
    O = "O",
    EMPTY = ""
}

enum BorderClass {
    TOP = "top-border",
    BOTTOM = "bottom-border",
    LEFT = "left-border",
    RIGHT = "right-border",
    SEPERATOR = " "
}

// TODO overvej om der skal lægge noget state ifht om den her square er x eller o i dette komponent i stedet for i Board.
// Board er ved at blive til et god-component
export const Square: React.FC<SquareProps> = (props) => {
    const getCharacter = () => {
        if (props.isX) return SquareCharacter.X;
        if (props.isO) return SquareCharacter.O;
        return SquareCharacter.EMPTY;
    }

    const getBorderClasses = () => {
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

        return borderClasses.join(BorderClass.SEPERATOR);
    }

    const character = getCharacter();
    const borderClasses = getBorderClasses();

    return <div onClick={() => props.chooseSquare(props.position)} className={"square" + BorderClass.SEPERATOR + borderClasses}>
        {character}
    </div>
}