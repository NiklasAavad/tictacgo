import { BackendPosition, Position } from "../utility/Position"

export const adaptPosition = (position: BackendPosition): Position => {
    return position + 1;
}