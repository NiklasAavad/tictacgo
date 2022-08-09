import React from "react";
import { EmptyString, SquareCharacter } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";

export type Result = {
    winningCombination: Position[],
    winningCharacter: SquareCharacter;
}

type SetBoardFunction = React.Dispatch<React.SetStateAction<(SquareCharacter | EmptyString)[]>>
type SetResultFunction = React.Dispatch<React.SetStateAction<Result | undefined>>
type SetIsGameOverFunction = React.Dispatch<React.SetStateAction<boolean>>
type SetIsGameStartedFunction = React.Dispatch<React.SetStateAction<boolean>>

export type GameContextMutator = {
    setBoard: SetBoardFunction,
    setResult: SetResultFunction,
    setIsGameOver: SetIsGameOverFunction,
    setIsGameStarted: SetIsGameStartedFunction,
}

// TODO hvis det er nødvendigt at sætte playerInTurn i GameContext, så skal den muligvis også med som parameter.
// I så fald vil det måske give mening at lave en type GameContextMutator, som har begge funktioner på sig.
export type GameService = (gameContextMutator: GameContextMutator) => {
    // reqired methods
    startGame: () => void,
    chooseSquare: (position: Position) => void,

    // optional methods, used for offline services.
    getResult?: () => Result | undefined,
    isGameOver?: () => boolean,
    isChoiceValid?: (position: Position) => boolean,
    changePlayerInTurn?: () => void
}