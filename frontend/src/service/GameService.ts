import React from "react";
import { SquareCharacter, SquareType } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";

export type Result = {
    winningCombination: Position[],
    winningCharacter: SquareCharacter;
}

type SetLatestSquareFunction = React.Dispatch<React.SetStateAction<SquareType | undefined>>
type SetResultFunction = React.Dispatch<React.SetStateAction<Result | undefined>>
type SetIsGameOverFunction = React.Dispatch<React.SetStateAction<boolean>>

export type GameContextMutator = {
    setLatestSquare: SetLatestSquareFunction,
    setResult: SetResultFunction,
    setIsGameOver: SetIsGameOverFunction,
}

// TODO hvis det er nødvendigt at sætte playerInTurn i GameContext, så skal den muligvis også med som parameter.
// I så fald vil det måske give mening at lave en type GameContextMutator, som har begge funktioner på sig.
export type GameService = (gameContextMutator: GameContextMutator) => {
    startGame: () => void,
    getResult: () => Result | undefined,
    isGameOver: () => boolean,
    isChoiceValid: (position: Position) => boolean,
    chooseSquare: (position: Position) => void,
    changePlayerInTurn: () => void
}