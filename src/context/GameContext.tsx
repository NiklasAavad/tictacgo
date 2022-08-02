import { createContext, PropsWithChildren, useContext, useEffect, useState } from 'react';
import { SquareType } from '../components/Content/Game/Square/Square';
import OfflineMultiplayerGameService from '../service/OfflineMultiplayerGameService';
import { Position } from '../utility/Position';

type GameContextType = {
    latestSquare: SquareType | null,
    winningCombination: Position[] | null,
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

export const GameProvider = ({ children }: PropsWithChildren) => {
    const [latestSquare, setLatestSquare] = useState<SquareType | null>(null);
    const [winningCombination, setWinningCombination] = useState<Position[] | null>(null)
    const [isGameStarted, setIsGameStarted] = useState(false);

    // Er det et problem at alle klienter spørger om spillet er slut? Ikke hvis man kan undgå at skulle broadcaste beskeden til alle.
    useEffect(() => {
        if (OfflineMultiplayerGameService.isGameOver()) {
            endGame();
        }
    }, [latestSquare])

    const chooseSquare = (position: Position) => {
        const latestSquare = OfflineMultiplayerGameService.chooseSquare(position);

        if (latestSquare) {
            setLatestSquare(latestSquare);
            OfflineMultiplayerGameService.changePlayerInTurn();
        }
    }

    const endGame = () => {
        const [winningCombination, winningCharacter] = OfflineMultiplayerGameService.getWinningCombination();

        if (winningCombination) {
            console.log(`Game was won by ${winningCharacter}!`)
            setWinningCombination(winningCombination)
        } else {
            console.log("Game was tied...");
        }

        const waitForNewGame = setTimeout(() => setIsGameStarted(false), 2500)
        return () => clearTimeout(waitForNewGame);
    }

    const startGame = () => {
        OfflineMultiplayerGameService.startGame();
        setLatestSquare(null);
        setWinningCombination(null);
        setIsGameStarted(true);
    }

    const exposedValues = {
        latestSquare,
        winningCombination,
        isGameStarted,
        startGame,
        chooseSquare
    }


    return <GameContext.Provider value={exposedValues}>
        {children}
    </GameContext.Provider>;
}