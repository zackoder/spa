import { signin, signout, signup } from "./signin.js";
import { setupPage, trackscroll, setupSPA, createHTMLel } from "./helpers.js";

import { notFound } from "./errpage.js";

export const socket = new WebSocket("ws://localhost:8080/ws");

socket.addEventListener("open", () => {
  console.log("connected..");
});

document.addEventListener("DOMContentLoaded", async () => {
  setupSPA();
});

export const routes = {
  "/signin": signin,
  "/signup": signup,
  "/signout": signout,
  "/404": notFound,
  "/": async () => {
    try {
      setupPage();
    } catch (err) {
      console.log(err);
      location.href = "/signin";
    }
  },
};

trackscroll();
