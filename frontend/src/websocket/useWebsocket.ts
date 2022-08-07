import { useCallback, useMemo } from "react";

export type MessageCallback = (msg: MessageEvent) => void;

const BASE_URL = "localhost:8080"

export const useWebsocket = (name: string | undefined) => {
    const userName = name || "Unknown User";
    const socket = useMemo(() => new WebSocket(`ws://${BASE_URL}/ws?name=${userName}`), [userName]);

    const connect = useCallback((messageCallback: MessageCallback) => {
        console.log("Attempting connection...");

        socket.onopen = () => {
            console.log("Successfully connected");
        };

        socket.onmessage = (msg: MessageEvent) => {
            console.log(msg);
            messageCallback(msg);
        };

        socket.onclose = (event: CloseEvent) => {
            console.log("Socket closed connection", event);
        };

        socket.onerror = (error: Event) => {
            console.log("Socket error:", error);
        };
    }, [socket]);

    const sendMessage = useCallback((msg: string) => {
        console.log("Sending message:", msg)
        socket.send(msg);
    }, [socket]);

    return { connect, sendMessage };
}