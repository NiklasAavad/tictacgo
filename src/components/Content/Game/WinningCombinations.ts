import { Position } from "./Position";

export const WINNING_COMBINATIONS = [
    // Horiziontal combinations
    [Position.TOP_LEFT, Position.TOP_CENTER, Position.TOP_RIGHT],
    [Position.CENTER_LEFT, Position.CENTER, Position.CENTER_RIGHT],
    [Position.BOTTOM_LEFT, Position.BOTTOM_CENTER, Position.BOTTOM_RIGHT],
    
    // Vertical combinations
    [Position.TOP_LEFT, Position.CENTER_LEFT, Position.BOTTOM_LEFT],
    [Position.TOP_CENTER, Position.CENTER, Position.BOTTOM_CENTER],
    [Position.TOP_RIGHT, Position.CENTER_RIGHT, Position.BOTTOM_RIGHT],

    // Diagonal combinations
    [Position.TOP_LEFT, Position.CENTER, Position.BOTTOM_RIGHT],
    [Position.BOTTOM_LEFT, Position.CENTER, Position.TOP_RIGHT]
];