import "./Board.css"
import { Row } from "./Row"

export const Board: React.FC = () => {
    return <div className="board">
        <Row/>
        <Row/>
        <Row/>
    </div>
}