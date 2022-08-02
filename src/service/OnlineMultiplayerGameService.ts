import { SquareType } from "../components/Content/Game/Square/Square";
import { Position } from "../utility/Position";
import { GameService, Result } from "./GameService";

const OnlineMultiplayerGameService: GameService = {
    startGame: function (): void {
        throw new Error("Function not implemented.");
    },
    getResult: function (): Result | undefined {
        throw new Error("Function not implemented.");
    },
    isGameOver: function (): boolean {
        throw new Error("Function not implemented.");
    },
    isChoiceValid: function (position: Position): boolean {
        throw new Error("Function not implemented.");
    },
    chooseSquare: function (position: Position): SquareType {
        throw new Error("Function not implemented.");
    },
    changePlayerInTurn: function (): void {
        throw new Error("Function not implemented.");
    }
}

export default OnlineMultiplayerGameService;