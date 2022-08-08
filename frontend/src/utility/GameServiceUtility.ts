import { EmptyString, SquareCharacter } from "../components/Content/Game/Square/Square";
import { Position } from "./Position";
import { WINNING_COMBINATIONS } from "./WinningCombinations";

export const hasPlayerWon = (playerPositions: Position[]) => {
    return WINNING_COMBINATIONS.some(combination => {
        return combination.every(position => playerPositions.includes(position))
    });
};

export function getRandomItem<Type>(list: Type[]): Type {
    return list[Math.floor((Math.random() * list.length))];
}

export const getEmptyBoard = (): (SquareCharacter | EmptyString)[] => {
    const EMPTY_BOARD: (SquareCharacter | EmptyString)[] = ["", "", "", "", "", "", "", "", ""];
    return EMPTY_BOARD;
}