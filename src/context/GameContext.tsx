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

    // TODO Det her kommer ikke til at virke online. Her skal latestSquare blive hentet på en anden måde (når spilleren skifter tur?).
    const chooseSquare = (position: Position) => {
        if (!OfflineMultiplayerGameService.isChoiceValid(position)) {
            return;
        }

        const latestSquare = OfflineMultiplayerGameService.chooseSquare(position);
        setLatestSquare(latestSquare);
        OfflineMultiplayerGameService.changePlayerInTurn();
    }

    const endGame = () => {
        const result = OfflineMultiplayerGameService.getResult();

        if (result) {
            console.log(`Game was won by ${result.winningCharacter}!`)
            setWinningCombination(result.winningCombination)
        } else {
            console.log("Game was tied...");
        }

        const waitForGameToEnd = setTimeout(() => setIsGameStarted(false), 2500)
        return () => clearTimeout(waitForGameToEnd);
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