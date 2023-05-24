import { GAMEINFO } from "../context/GameContext";
import { Board, GameContextMutator } from "../service/GameService";
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

export const getEmptyBoard = (): Board => {
    const EMPTY_BOARD: Board = ["", "", "", "", "", "", "", "", ""];
    return EMPTY_BOARD;
}

export const startGameForMutator = (gameContextMutator: GameContextMutator) => {
	gameContextMutator.setBoard(getEmptyBoard());
	gameContextMutator.setLatestGameInfoMessage(GAMEINFO.NEW_GAME_STARTED);
	gameContextMutator.setIsGameStarted(true);
	gameContextMutator.setResult(undefined);
	gameContextMutator.setIsGameOver(false);
}
