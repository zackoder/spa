import { addevents } from "./helpers.js";
import { signin, signup } from "./signin.js";

const root = document.querySelector(".root");

document.addEventListener("DOMContentLoaded", () => {
  let path = location.pathname;
  if (path === "/signin") {
    signin();
  } else if (path === "/signup") {
    signup();
  }
});


