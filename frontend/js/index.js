import { signin, signout, signup } from "./signin.js";
import { chatbox } from "./chatFunctionality.js";
import { navbar, searchBar } from "./navbar.js";
import {
  addPostPopUp,
  createHTMLel,
  fetchcategory,
  handlescroll,
  showPosts,
} from "./helpers.js";
import { root } from "./navbar.js";
import { notFound } from "./errpage.js";
document.addEventListener("DOMContentLoaded", async () => {
  setupSPA();
});

const routes = {
  "/signin": signin,
  "/signup": signup,
  "/signout": signout,
  "/404": notFound,
  "/": async () => {
    try {
      setupPage();
      showPosts(location.pathname);
    } catch (err) {
      console.log(err);
      location.href = "/signin";
    }
    // document.querySelector("#startChat")?.addEventListener("click", chatbox);
  },
};

let scroll = window.scrollY;

document.addEventListener("scrollend", () => {
  console.log(window.innerHeight + window.scrollY);
  if (window.innerHeight + window.scrollY >= root.offsetHeight - 200) {
    console.log("hi");
    handlescroll(scroll);
  }
});

async function setupPage() {
  if (!document.querySelector(".header")) {
    await navbar();
    searchBar();
    const style = createHTMLel("link", "", "", {
      key: "href",
      value: "/frontend/style/post.css",
    });

    style.rel = "stylesheet";
    const title = createHTMLel("title", "", "Forum");
    document.head.append(style, title);
    addPostPopUp();
  }
  // showPosts(location.pathname);
}
export let offset = 0;

function setupSPA() {
  root.addEventListener("click", (e) => {
    const link = e.target.closest("a");
    if (link && link.origin === location.origin) {
      e.preventDefault();
      offset = 0;
      document.querySelector(".postscontainer").innerHTML = "";
      navigateTo(link.pathname);
    }
  });

  window.addEventListener("popstate", () => handleRoute(location.pathname));

  handleRoute(location.pathname);
}

export function getOffset() {
  return offset;
}

export function setOffset(value) {
  offset = value;
}

function navigateTo(path) {
  history.pushState({}, "", path);
  handleRoute(path);
}

// Function to handle route changes
async function handleRoute(path) {
  if (routes[path]) {
    await routes[path]();
  } else if (isUsernamePath(path)) {
    await fetchUserProfile(path.substring(1));
    setupPage();
  } else if (path.startsWith("/category/")) {
    await fetchcategory(path);
    setupPage();
  } else {
    navigateTo("/404");
  }
}

function isUsernamePath(path) {
  return /^\/[a-zA-Z0-9_-]+$/.test(path);
}

async function fetchUserProfile(username) {
  try {
    let response = await fetch(`/${username}`);
    if (!response.ok) throw new Error("User not found");

    let userData = await response.json();
    console.log("User Profile:", userData);
    renderUserProfile(userData);
  } catch (error) {
    console.error(error);
    navigateTo("/404");
  }
}
