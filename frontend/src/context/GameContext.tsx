import { createContext, PropsWithChildren, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { SquareCharacter } from '../components/Content/Game/Square/Square';
import { Board, GameService, Result } from '../service/GameService';
import { getEmptyBoard } from '../utility/GameServiceUtility';
import { Position } from '../utility/Position';

type GameContextType = {
    board: Board,
    latestGameInfoMessage: GameInfoMessage,
    result: Result | undefined,
    isGameStarted: boolean,
    chooseSquare: (position: Position) => void,
    startGame: () => void,
    selectCharacter: (character: SquareCharacter) => void
}

const GameContext = createContext<GameContextType | undefined>(undefined);

export const useGameContext = () => {
    const context = useContext(GameContext);

    if (context === undefined) {
        throw new Error("useGameContext must be used within a GameProvider!")
    }

    return context;
}

export const GAMEINFO = {
	X_WON: "The game has been won by X!",
    O_WON: "The game has been won by O!",
    TIE: "The game has been tied...",
    START_NEW_GAME: "Click on the 'New Game' button to begin!",
    NEW_GAME_STARTED: "A new game has just begun! X begins.",
    X_SELECTED: "X has been selected",
    O_SELECTED: "O has been selected",
} as const;

export type GameInfoMessage = typeof GAMEINFO[keyof typeof GAMEINFO];

type GameProviderProps = {
    gameServiceProvider: GameService
}

const TIMEOUT_PERIOD = 2500; // ms!

export const GameProvider = ({ gameServiceProvider, children }: PropsWithChildren<GameProviderProps>) => {
    const [board, setBoard] = useState<Board>(getEmptyBoard());
    const [latestGameInfoMessage, setLatestGameInfoMessage] = useState<GameInfoMessage>(GAMEINFO.START_NEW_GAME);
    const [isGameStarted, setIsGameStarted] = useState(false);
    const [result, setResult] = useState<Result | undefined>(undefined);
    const [isGameOver, setIsGameOver] = useState<boolean>(false);

    const gameContextMutator = useMemo(() => {
        return { setBoard, setResult, setIsGameOver, setIsGameStarted, setLatestGameInfoMessage }
    }, []);

    const gameService = useMemo(() => gameServiceProvider(gameContextMutator), [gameServiceProvider, gameContextMutator]);

    const startGame = () => {
        gameService.startGame();
    }

    const chooseSquare = (position: Position) => {
        gameService.chooseSquare(position);
    }

    const selectCharacter = (character: SquareCharacter) => {
        gameService.selectCharacter(character);
    }

    const getWinningMessage = useCallback((result: Result): GameInfoMessage => {
        const xWon = result.winningCharacter === SquareCharacter.X;
        if (xWon) {
            return GAMEINFO.X_WON;
        }
        else {
            return GAMEINFO.O_WON;
        }
    }, []);

    const endGame = useCallback(() => {
        if (result) {
            const newGameMessage = getWinningMessage(result);
            setLatestGameInfoMessage(newGameMessage);
        } else {
            setLatestGameInfoMessage(GAMEINFO.TIE);
        }

        const waitForGameToEnd = setTimeout(() => {
            setLatestGameInfoMessage(GAMEINFO.START_NEW_GAME);
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
        result,
        isGameStarted,
        startGame,
        chooseSquare,
        selectCharacter
    }

    return <GameContext.Provider value={exposedValues}>
        {children}
    </GameContext.Provider>;
}
