import { creatcategories, createHTMLel, sendPost } from "./helpers.js";

export const root = document.getElementById("root");

export const navbar = async (nickname) => {
  const style = createHTMLel("link", "", "", {
    key: "href",
    value: "/frontend/style/navbar.css",
  });
  style.rel = "stylesheet";
  document.head.appendChild(style);

  const header = createHTMLel("header", "header");
  const logo = createHTMLel("a", "logo", "forum", { key: "href", value: "/" });

  const profile = createHTMLel("a", "link", "Profile", {
    key: "href",
    value: `/${nickname}`,
  });

  const logout = createHTMLel("a", "link", "Sing Out", {
    key: "href",
    value: "/signout",
  });
  const ul = createHTMLel("ul", "navbarUl");
  ul.append(profile, logout);
  header.append(logo, ul);
  root.appendChild(header);
};

export const searchBar = () => {
  const style = createHTMLel("link", "", "", {
    key: "href",
    value: "/frontend/style/index.css",
  });

  style.rel = "stylesheet";
  document.head.appendChild(style);

  const div = createHTMLel("div", "searchbar expanded");
  const arrow = createHTMLel("span", "arrow rotated");
  arrow.innerHTML = "&#9660;";
  arrow.addEventListener("click", () => {
    div.classList.toggle("expanded");
    arrow.classList.toggle("rotated");
  });

  let lastScrollY = window.scrollY;

  window.addEventListener("scroll", () => {
    const currentScrollY = window.scrollY;
    if (!currentScrollY) {
      div.classList.add("expanded");
      arrow.classList.add("rotated");
    } else {
      div.classList.remove("expanded");
      arrow.classList.remove("rotated");
    }
    if (currentScrollY > lastScrollY) {
      div.classList.remove("expanded");
      arrow.classList.remove("rotated");
    }

    lastScrollY = currentScrollY;
  });

  const categoriesSilder = createHTMLel("div", "category-slider");

  creatcategories(categoriesSilder, "a");

  const addposticon = createHTMLel("img", "addpost", "", {
    key: "src",
    value: "/frontend/images/add.png",
  });

  addposticon.addEventListener("click", () => {
    const postcontainer = document.querySelector(".addPostContainer");
    postcontainer.classList.toggle("show");
  });
  div.append(arrow, categoriesSilder, addposticon);
  root.appendChild(div);
};
