import { createContext, PropsWithChildren, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { EmptyString, SquareCharacter, SquareType } from '../components/Content/Game/Square/Square';
import { GameService, Result } from '../service/GameService';
import { EMPTY_BOARD } from '../utility/GameServiceUtility';
import { Position } from '../utility/Position';

type GameContextType = {
    board: (SquareCharacter | EmptyString)[],
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
    NEW_GAME_STARTED = "A new game has just begun! X begins."
}

type GameProviderProps = {
    gameServiceProvider: GameService
}

const TIMEOUT_PERIOD = 2500; // ms!

export const GameProvider = ({ gameServiceProvider, children }: PropsWithChildren<GameProviderProps>) => {
    const [board, setBoard] = useState<(SquareCharacter | EmptyString)[]>(EMPTY_BOARD);
    const [latestGameInfoMessage, setLatestGameInfoMessage] = useState<GameInfoMessage>(GameInfoMessage.START_NEW_GAME);
    const [winningCombination, setWinningCombination] = useState<Position[] | undefined>(undefined)
    const [isGameStarted, setIsGameStarted] = useState(false);
    const [result, setResult] = useState<Result | undefined>(undefined);
    const [isGameOver, setIsGameOver] = useState<boolean>(false);
    const [playerInTurn, setPlayerInTurn] = useState<SquareCharacter>(SquareCharacter.X);

    const gameContextMutator = useMemo(() => {
        return { setBoard, setResult, setIsGameOver }
    }, [setBoard, setResult, setIsGameOver]);

    const gameService = useMemo(() => gameServiceProvider(gameContextMutator), [gameServiceProvider, gameContextMutator]);

    const startGame = () => {
        gameService.startGame();
        setBoard(EMPTY_BOARD);
        setLatestGameInfoMessage(GameInfoMessage.NEW_GAME_STARTED);
        setWinningCombination(undefined);
        setIsGameStarted(true);
        setResult(undefined);
        setIsGameOver(false);
    }

    const chooseSquare = (position: Position) => {
        gameService.chooseSquare(position);
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
    }, [getWinningMessage, result]);

    useEffect(() => {
        if (isGameOver) {
            endGame();
        }
    }, [isGameOver, endGame])

    const exposedValues = {
        board,
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