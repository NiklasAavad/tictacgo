import { useCallback, useMemo } from "react";
import { BASE_URL } from "../api/Api";

export type MessageCallback = (msg: MessageEvent) => void;


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

    const sendChatMessage = useCallback((msg: string) => {
        console.log("Sending chat message:", msg)

        const jsonMsg = JSON.stringify({
            type: 0,
            sender: userName,
            body: msg
        })
        
        socket.send(jsonMsg);
    }, [socket, userName]);

    return { connect, sendChatMessage };
}