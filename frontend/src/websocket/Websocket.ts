const socket = new WebSocket("ws://localhost:8080/ws");

export const connect = () => {
    console.log("Attempting connection...");

    socket.onopen = () => {
        console.log("Successfully connected");
    };

    socket.onmessage = (msg: MessageEvent) => {
        console.log(msg);
    };

    socket.onclose = (event: CloseEvent) => {
        console.log("Socket closed connection", event);
    };

    socket.onerror = (error: Event) => {
        console.log("Socket error:", error);
    };
};

// TODO fix any
export const sendMsg = (msg: any) => {
    console.log("Sending message:", msg)
    socket.send(msg);
}