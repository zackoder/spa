import { createHTMLel } from "./helpers.js";

export const chatbox = () => {
  console.log("chatbox function executed!");
  let socket = new WebSocket("ws://localhost:8080/ws");
  socket.addEventListener("open", (event) => {
    console.log("web socket connection");
    socket.send("hello world");
  });
  socket.addEventListener("error", (event) => {
    console.error("WebSocket error:", event);
  });

  // socket.onopen = (event) => {
  //     console.log("web socket connection")
  //     socket.send("hello world")
  // }
  socket.onmessage = (event) => {
    const message = createHTMLel("div", "message", event.data);
  };
};
