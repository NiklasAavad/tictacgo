import { getRandomItem } from "../utility/GameServiceUtility";
import { Position, positions } from "../utility/Position";
import { GameContextMutator, GameService, Result } from "./GameService";
import OfflineMultiplayerGameService from "./OfflineMultiplayerGameService";

const OfflineSingleplayerGameService: GameService = (gameContextMutator: GameContextMutator) => {
    const decoratee = OfflineMultiplayerGameService(gameContextMutator);

    const choosePositionByBot = () => {
        let position = undefined;
        do {
            position = getRandomItem(positions);
        } while (position === undefined || !isChoiceValid(position));

        return position;
    }

    const chooseSquareByBot = () => {
        const positionChosenByBot = choosePositionByBot();
        decoratee.chooseSquare(positionChosenByBot);
        decoratee.changePlayerInTurn();
    }

    const startGame = (): void => {
        decoratee.startGame();
    };

    const getResult = (): Result | undefined => {
        return decoratee.getResult();
    };

    const isGameOver = (): boolean => {
        return decoratee.isGameOver();
    };

    const isChoiceValid = (position: Position): boolean => {
        return decoratee.isChoiceValid(position);
    };

    const chooseSquare = (position: Position): void => {
        decoratee.chooseSquare(position);

        if (isGameOver()) {
            return;
        }

        setTimeout(() => chooseSquareByBot(), 500);
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