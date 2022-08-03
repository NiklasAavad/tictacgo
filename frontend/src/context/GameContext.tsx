import { createContext, PropsWithChildren, useCallback, useContext, useEffect, useState } from 'react';
import { SquareCharacter, SquareType } from '../components/Content/Game/Square/Square';
import { GameService, Result } from '../service/GameService';
import { Position } from '../utility/Position';

type GameContextType = {
    latestSquare: SquareType | undefined,
    latestGameInfoMessage: GameInfoMessage,
    winningCombination: Position[] | undefined,
    isGameStarted: boolean,
    chooseSquare: (position: Position) => void,
    startGame: () => void
}

const GameContext = createContext<GameContextType | undefined>(undefined);
export const useGameContext = () => {
    const context = useContext(GameContext);

    if (context === undefined) {
        throw new Error("useGameContext must be used within a GameProvider!")
    }

    return context;
}

enum GameInfoMessage {
    X_WON = "The game has been won by X!",
    O_WON = "The game has been won by O!",
    TIE = "The game has been tied...",
    START_NEW_GAME = "Click on the 'New Game' button to begin!",
    NEW_GAME_STARTED = "A new game has just begun! Good luck."
}

type GameProviderProps = {
    gameService: GameService
}

const TIMEOUT_PERIOD = 2500; // ms!

export const GameProvider = ({ gameService, children }: PropsWithChildren<GameProviderProps>) => {
    const [latestSquare, setLatestSquare] = useState<SquareType | undefined>(undefined);
    const [latestGameInfoMessage, setLatestGameInfoMessage] = useState<GameInfoMessage>(GameInfoMessage.START_NEW_GAME);
    const [winningCombination, setWinningCombination] = useState<Position[] | undefined>(undefined)
    const [isGameStarted, setIsGameStarted] = useState(false);

    const startGame = () => {
        gameService.startGame();
        setLatestSquare(undefined);
        setLatestGameInfoMessage(GameInfoMessage.NEW_GAME_STARTED);
        setWinningCombination(undefined);
        setIsGameStarted(true);
    }

    // TODO Det her kommer ikke til at virke online. Her skal latestSquare blive hentet på en anden måde (når spilleren skifter tur?).
    const chooseSquare = (position: Position) => {
        if (!gameService.isChoiceValid(position)) {
            return;
        }

        const latestSquare = gameService.chooseSquare(position);
        setLatestSquare(latestSquare);
        gameService.changePlayerInTurn();
    }

    const getWinningMessage = useCallback((result: Result): GameInfoMessage => {
        const xWon = result.winningCharacter === SquareCharacter.X;
        if (xWon) {
            return GameInfoMessage.X_WON;
        }
        else {
            return GameInfoMessage.O_WON;
        }
    }, []);

    const endGame = useCallback(() => {
        const result = gameService.getResult();

        if (result) {
            const newGameMessage = getWinningMessage(result);
            setLatestGameInfoMessage(newGameMessage);
            setWinningCombination(result.winningCombination)
        } else {
            setLatestGameInfoMessage(GameInfoMessage.TIE);
        }

        const waitForGameToEnd = setTimeout(() => {
            setLatestGameInfoMessage(GameInfoMessage.START_NEW_GAME);
            setIsGameStarted(false)
        }, TIMEOUT_PERIOD);

        return () => clearTimeout(waitForGameToEnd);
    }, [gameService, getWinningMessage]);

    // Er det et problem at alle klienter spørger om spillet er slut? Ikke hvis man kan undgå at skulle broadcaste beskeden til alle.
    useEffect(() => {
        if (gameService.isGameOver()) {
            endGame();
        }
    }, [latestSquare, gameService, endGame])

    const exposedValues = {
        latestSquare,
        latestGameInfoMessage,
        winningCombination,
        isGameStarted,
        startGame,
        chooseSquare
    }

    return <GameContext.Provider value={exposedValues}>
        {children}
    </GameContext.Provider>;
}