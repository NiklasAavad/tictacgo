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

    // TODO der er en bug ifht Position.TOP_LEFT. Det er fixet hvis TOP_LEFt starter fra 1, men så passer det ikke med board og backend.
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
        if (!decoratee.getResult) {
            throw new Error("getResult could not be delegated to decoratee.")
        }

        return decoratee.getResult();
    };

    const isGameOver = (): boolean => {
        if (!decoratee.isGameOver) {
            throw new Error("isGameOver could not be delegated to decoratee.")
        }

        return decoratee.isGameOver();
    };

    const isChoiceValid = (position: Position): boolean => {
        if (!decoratee.isChoiceValid) {
            throw new Error("isChoiceValid could not be delegated to decoratee.")
        }

        if (isBotChoosing) {
            return false;
        }

        return decoratee.isChoiceValid(position);
    };

    const chooseSquare = (position: Position): void => {
        if (!isChoiceValid(position)) {
            return;
        }

        decoratee.chooseSquare(position);
        addSquareToBoard(SquareCharacter.X, position)

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
        if (!decoratee.changePlayerInTurn) {
            throw new Error("changePlayerInTurn could not be delegated to decoratee.")
        }

        decoratee.changePlayerInTurn();
    };

    const selectCharacter = (character: SquareCharacter): void => {
        decoratee.selectCharacter(character);
    }

    return {
        startGame,
        getResult,
        isGameOver,
        isChoiceValid,
        chooseSquare,
        changePlayerInTurn,
        selectCharacter
    }
}

export default OfflineSingleplayerGameService;