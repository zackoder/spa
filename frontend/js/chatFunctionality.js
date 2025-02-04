import { createHTMLel } from "./helpers.js";

export const chatbox = () => {
    console.log("hello");

    let socket = new WebSocket("ws://localhost:8088/ws");
    socket.addEventListener("open", (event) => {
        console.log("web socket connection")
        socket.send("hello world")
    })

    // socket.onopen = (event) => {
    //     console.log("web socket connection")
    //     socket.send("hello world")
    // }
    socket.onmessage = (event) => {
        const message = createHTMLel("div", "message", event.data)
    }
}