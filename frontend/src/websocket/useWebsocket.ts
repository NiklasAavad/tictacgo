import { useMemo } from "react";

export type MessageCallback = (msg: MessageEvent) => void;

export const useWebsocket = () => {
    const socket = useMemo(() => new WebSocket("ws://localhost:8080/ws?name=Niklas"), []);

    const connect = (messageCallback: MessageCallback) => {
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
    };

    const sendMessage = (msg: string) => {
        console.log("Sending message:", msg)
        socket.send(msg);
    };

    return { connect, sendMessage };
}