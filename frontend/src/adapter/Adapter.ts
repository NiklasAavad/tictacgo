import { EmptyString, SquareCharacter } from "../components/Content/Game/Square/Square";
import { Board } from "../service/GameService";

export type BackendWinningCharacter = 0 | 1;
export type BackendSquareChacter = 0 | 1 | 2;

// TODO Bliver inputtet valideret inden vi når her ind? Eller skal vi validere noget her?

export const adaptBoard = (backendBoard: BackendSquareChacter[]): Board => {
    return backendBoard.map(character => adaptCharacter(character))
}

export const adaptWinningCharacter = (backendWinningCharacter: BackendWinningCharacter): SquareCharacter => {
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