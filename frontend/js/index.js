import { signin, signout, signup } from "./signin.js";
import { chatbox } from "./chatFunctionality.js";
import { navbar, searchBar } from "./navbar.js";
import {
  addPostPopUp,
  createHTMLel,
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
      let res = await fetch("/getNickName");

      if (!res.ok) {
        location.href = "/signin";
      } else {
        let data = await res.json();

        let nickname = data.nickname;
        navbar(nickname);
        searchBar();
        const style = createHTMLel("link", "", "", {
          key: "href",
          value: "/frontend/style/post.css",
        });

        style.rel = "stylesheet";
        const title = createHTMLel("title", "", "Forum");
        document.head.append(style, title);
        addPostPopUp();
        showPosts("/posts");

        let scroll = window.scrollY;
        document.addEventListener("scrollend", () => {
          console.log(window.innerHeight + window.scrollY);
          console.log(root.offsetHeight);
          if (window.innerHeight + window.scrollY >= root.offsetHeight - 200) {
            handlescroll(scroll);
          }
        });
      }
    } catch (err) {
      console.log(err);
      location.href = "/signin";
    }
    document.querySelector("#startChat")?.addEventListener("click", chatbox);
  },
};

function setupSPA() {
  root.addEventListener("click", (e) => {
    const link = e.target.closest("a");
    if (link && link.origin === location.origin) {
      e.preventDefault();
      navigateTo(link.pathname);
    }
  });

  window.addEventListener("popstate", () => handleRoute(location.pathname));

  handleRoute(location.pathname);
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
  } else if (path.includes("/category/")) {
    console.log(path);
    await fetchcategory(path);
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
