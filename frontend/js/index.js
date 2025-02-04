import { signin, signup } from "./signin.js";
import { chatbox } from "./chatFunctionality.js";

const root = document.querySelector(".root");

document.addEventListener("DOMContentLoaded", () => {
  let path = location.pathname;
  if (path === "/signin") {
    signin();
  } else if (path === "/signup") {
    signup();
  } else if (path === "/ws") {
    console.log("/ws");
    chatbox()
  }
});


