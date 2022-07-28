import "./Row.css"
import { Square } from "./Square"

export const Row: React.FC = () => {
    return <div className="row">
        <Square/>
        <Square/>
        <Square/>
    </div>
}