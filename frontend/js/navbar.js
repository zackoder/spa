import { creatcategories, createHTMLel, layout } from "./helpers.js";
export let user = "";

export const root = document.getElementById("root");
root.scrollTop = root.scrollHeight;

export const navbar = async () => {
  let res = await fetch("/getNickName");
  if (!res.ok) {
    location.href = "/signin";
    return;
  } else {
    let data = await res.json();
    const style = createHTMLel("link", "", "", {
      key: "href",
      value: "/frontend/style/navbar.css",
    });
    style.rel = "stylesheet";
    document.head.appendChild(style);

    const header = createHTMLel("header", "header");
    const logo = createHTMLel("a", "logo", "forum", {
      key: "href",
      value: "/",
    });

    const profile = createHTMLel("a", "link", "Profile", {
      key: "href",
      value: `/${data.nickname}`,
    });

    user = data.nickname;
    console.log(data.nickname);

    const logout = createHTMLel("a", "link", "Sing Out", {
      key: "href",
      value: "/signout",
    });
    const ul = createHTMLel("ul", "navbarUl");
    ul.append(profile, logout);
    header.append(logo, ul);
    root.appendChild(header);
  }
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
    const usersContainer = document.querySelector(".usersContainer");
    if (!currentScrollY) {
      div.classList.add("expanded");
      arrow.classList.add("rotated");
      usersContainer?.classList.remove("up");
    } else {
      div.classList.remove("expanded");
      arrow.classList.remove("rotated");
      usersContainer?.classList.add("up");
    }
    if (currentScrollY > lastScrollY) {
      div.classList.remove("expanded");
      arrow.classList.remove("rotated");
      usersContainer?.classList.add("up");
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
    layout.classList.toggle("layout");
    if (layout.classList.contains("layout"))
      document.body.style.overflow = "hidden";
  });
  div.append(arrow, categoriesSilder, addposticon);
  root.appendChild(div);
};
