import { SquareCharacter } from "../components/Content/Game/Square/Square";
import { getRandomItem, hasPlayerWon } from "../utility/GameServiceUtility";
import { allPositions, Position } from "../utility/Position";
import { GameContextMutator, GameService, Result } from "./GameService";
import OfflineMultiplayerGameService from "./OfflineMultiplayerGameService";

const OfflineSingleplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    let x: Position[] = [];
    let o: Position[] = [];
    let freePositions = allPositions;
    const decoratee = OfflineMultiplayerGameService(gameContextMutator);
    let isBotChoosing = false;

    const addSquareToBoard = (playerInTurn: SquareCharacter, position: Position) => {
        if (playerInTurn === SquareCharacter.X) {
            x.push(position);
        } else {
            o.push(position);
        }
        freePositions = freePositions.filter(freePosition => freePosition !== position);
    }

    const choosePositionIfGameWouldEnd = (playerPositions: Position[]): Position | undefined => {
        for (const freePosition of freePositions) {
            const positionsIncludingFreePosition = [...playerPositions, freePosition];
            if (hasPlayerWon(positionsIncludingFreePosition)) {
                return freePosition;
            }
        }
        return undefined;
    }

    const choosePositionByBot = () => {
        const positionIfBotWouldWin = choosePositionIfGameWouldEnd(o);
        if (positionIfBotWouldWin) {
            return positionIfBotWouldWin;
        }

        const positionIfPlayerWouldWin = choosePositionIfGameWouldEnd(x);
        if (positionIfPlayerWouldWin) {
            return positionIfPlayerWouldWin;
        }

        if (freePositions.includes(Position.CENTER)) {
            return Position.CENTER;
        }

        return getRandomItem(freePositions);
    }

    const chooseSquareByBot = () => {
        const positionChosenByBot = choosePositionByBot();
        decoratee.chooseSquare(positionChosenByBot);
        addSquareToBoard(SquareCharacter.O, positionChosenByBot);
    }

    const startGame = (): void => {
        freePositions = allPositions;
        x = [];
        o = [];
        decoratee.startGame();
    };

    const getResult = (): Result | undefined => {
        return decoratee.getResult();
    };

    const isGameOver = (): boolean => {
        return decoratee.isGameOver();
    };

    const isChoiceValid = (position: Position): boolean => {
        if (isBotChoosing) {
            return false;
        }

        return decoratee.isChoiceValid(position);
    };

    const chooseSquare = (position: Position): void => {
        decoratee.chooseSquare(position);
        addSquareToBoard(SquareCharacter.X, position);

        if (isGameOver()) {
            return;
        }

        isBotChoosing = true;
        setTimeout(() => {
            chooseSquareByBot()
            isBotChoosing = false;
        }, 500);
    };

    const changePlayerInTurn = (): void => {
        decoratee.changePlayerInTurn();
    };

    return {
        startGame,
        getResult,
        isGameOver,
        isChoiceValid,
        chooseSquare,
        changePlayerInTurn
    }
}

export default OfflineSingleplayerGameService;