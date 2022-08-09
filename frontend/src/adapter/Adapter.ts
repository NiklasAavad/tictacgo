import { EmptyString, SquareCharacter } from "../components/Content/Game/Square/Square";

export type BackendWinningCharacter = 0 | 1;
export type BackendSquareChacter = 0 | 1 | 2;

export const adaptBoard = (backendBoard: BackendSquareChacter[]): (SquareCharacter | EmptyString)[] => {
    return backendBoard.map(square => adaptCharacter(square))
}

export const adaptWinningCharacter = (backendWinningCharacter: BackendWinningCharacter): SquareCharacter => {
    console.log("adapting character:", backendWinningCharacter);
    return backendWinningCharacter === 0 ? SquareCharacter.X : SquareCharacter.O;
}

const adaptCharacter = (backendCharacter: BackendSquareChacter): SquareCharacter | EmptyString => {
    switch (backendCharacter) {
        case 0:
            return SquareCharacter.X;
        case 1:
            return SquareCharacter.O;
        case 2:
            return "";
    }
}