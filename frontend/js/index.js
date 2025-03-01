import { signin, signout, signup } from "./signin.js";
import { setupPage, trackscroll, setupSPA } from "./helpers.js";

import { notFound } from "./errpage.js";

export const originalHTML = document.documentElement.innerHTML;

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
