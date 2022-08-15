import { useCallback, useMemo } from "react";
import { CHAT_WS_URL } from './../api/BackendApi';

export type MessageCallback = (msg: MessageEvent) => void;

const RECONNECT_INTERVAL = 1000; // ms!

export const useWebsocket = (name: string | undefined, messageCallback: MessageCallback) => {
    const userName = name || "Unknown User";
    const socket = useMemo(() => new WebSocket(`${CHAT_WS_URL}?name=${userName}`), [userName]);

    const connect = useCallback(() => {
        console.log("Attempting connection...");

        socket.onmessage = (msg: MessageEvent) => {
            console.log(msg);
            messageCallback(msg);
        };

        socket.onclose = (event: CloseEvent) => {
            console.log("Socket closed connection", event);
            console.log("Trying to reconnect...");
            setTimeout(() => connect(), RECONNECT_INTERVAL); // TODO Websocket is already in CLOSING or CLOSED state.
        };

        socket.onerror = (error: Event) => {
            console.log("Socket error:", error);
            console.log("Closing connection...");
            socket.close();
        };
    }, [socket, messageCallback]);

    const sendChatMessage = useCallback((msg: string) => {
        console.log("Sending chat message:", msg)
        socket.send(msg);
    }, [socket]);

    return { connect, sendChatMessage };
}