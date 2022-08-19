import React from "react";
import { EmptyString, SquareCharacter } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";

export type Result = {
    winningCombination: Position[],
    winningCharacter: SquareCharacter;
}

export type Board = (SquareCharacter | EmptyString)[];

type SetBoardFunction = React.Dispatch<React.SetStateAction<Board>>
type SetResultFunction = React.Dispatch<React.SetStateAction<Result | undefined>>
type SetIsGameOverFunction = React.Dispatch<React.SetStateAction<boolean>>
type SetIsGameStartedFunction = React.Dispatch<React.SetStateAction<boolean>>

export type GameContextMutator = {
    setBoard: SetBoardFunction,
    setResult: SetResultFunction,
    setIsGameOver: SetIsGameOverFunction,
    setIsGameStarted: SetIsGameStartedFunction,
}

export type GameService = (gameContextMutator: GameContextMutator) => {
    // required methods
    startGame: () => void,
    chooseSquare: (position: Position) => void,
    selectCharacter: (character: SquareCharacter) => void,

    // optional methods, used for offline services.
    getResult?: () => Result | undefined,
    isGameOver?: () => boolean,
    isChoiceValid?: (position: Position) => boolean,
    changePlayerInTurn?: () => void
}