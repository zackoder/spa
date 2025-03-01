import { navbar, searchBar } from "./navbar.js";
import { root } from "./navbar.js";

import { routes } from "./index.js";
import { getuser } from "./getusers.js";
import { socketEvents, upgradeconnection } from "./chatFunctionality.js";

export const createHTMLel = (
  name,
  Class,
  content = "",
  atrebute = { key: "", value: "" }
) => {
  let element = document.createElement(name);
  if (content == 0 || content) element.textContent = content;
  if (Class) element.className = Class;
  if (atrebute.key) element.setAttribute(atrebute.key, atrebute.value);
  return element;
};

export const layout = createHTMLel("div");
layout.addEventListener("click", () => {
  const post = document.querySelector(".show");
  document.body.style.overflow = "";
  layout.classList.toggle("layout");
  post.classList.toggle("show");
});

root.append(layout);
export function navigateTo(path) {
  history.pushState({}, "", path);
  handleRoute(path);
}

async function handleRoute(path) {
  if (routes[path]) {
    await routes[path]();
  } else if (path.startsWith("/category/") || isUsernamePath(path)) {
    setupPage();
  } else {
    navigateTo("/404");
  }
}

function isUsernamePath(path) {
  return /^\/[a-zA-Z0-9_-]+$/.test(path);
}

const main = createHTMLel("main", "main");
const sidebarLeft = createHTMLel("aside", "sidebar left-sidebar");
const sidebarRight = createHTMLel("aside", "sidebar right-sidebar");

export const fetchData = async (path, data) => {
  let resp = await fetch(path, {
    method: "POST",
    headers: {
      "Content-type": "application/json",
    },
    body: JSON.stringify(data),
  });
  return resp;
};

export const sendPost = async (title, content, categories, errp) => {
  if (!title) {
    errp.textContent = "Title can not be empty";
    errp.style.display = "block";
    return;
  }

  if (!content) {
    errp.textContent = "Content can not be empty";
    errp.style.display = "block";
    return;
  }

  if (categories.length === 0) {
    errp.textContent = "You need to choose at least one category.";
    errp.style.display = "block";
    return;
  }

  const data = {
    title: title,
    content: content,
    categories: categories,
  };

  const res = await fetchData("/addpost", data);
  if (!res.ok) {
    console.log("while adding a post the res is not ok ", res);
    return;
  }
  const newpost = await res.json();

  let path = decodeURIComponent(location.pathname);
  for (let i = 0; i < categories.length; i++) {
    if (path === "/" || path.endsWith(categories[i])) {
      offset++;
      creatPosts(
        document.querySelector(".postscontainer"),
        [newpost],
        "prepend"
      );
      return;
    }
  }
};

export const addPostPopUp = async () => {
  const div = createHTMLel("div", "addPostContainer");
  const h1 = createHTMLel("h1", "addPostHead", "Creat Post");
  const titleLbl = createHTMLel("label", "lbl", "Title", {
    key: "for",
    value: "titelInpt",
  });
  const titleinpt = createHTMLel("input", "inpt", "", {
    key: "placeholder",
    value: "Enter your title",
  });
  const contentLbl = createHTMLel("label", "lbl", "add content", {
    key: "for",
    value: "contentT",
  });
  const contnetinpt = createHTMLel("textarea", "contentErea", "", {
    key: "placeholder",
    value: "enter the content",
  });
  const categories = createHTMLel("div", "categories");
  await creatcategories(categories, "div");

  let chosenCategoreis = Array.from(categories.querySelectorAll(".category"));

  chosenCategoreis.forEach((btn) => {
    btn.addEventListener("click", () => {
      btn.classList.toggle("chosen");
    });
  });

  const submitbtn = createHTMLel("button", "submitbtn", "submit");
  const errp = createHTMLel("p", "errorp");

  submitbtn.addEventListener("click", () => {
    let chosenCate = [];
    chosenCategoreis.forEach((category) => {
      if (category.classList.contains("chosen")) {
        chosenCate.push(category.textContent);
      }
    });
    sendPost(titleinpt.value, contnetinpt.value, chosenCate, errp);
  });

  div.append(
    h1,
    titleLbl,
    titleinpt,
    contentLbl,
    contnetinpt,
    categories,
    errp,
    submitbtn
  );
  root.appendChild(div);
};

export const creatcategories = async (categoriesSlider, type) => {
  const left_arrow = createHTMLel("button", "arrows");
  left_arrow.innerHTML = "&lt;";
  const rgth_arrow = createHTMLel("button", "arrows");
  rgth_arrow.innerHTML = "&gt;";
  left_arrow.addEventListener("click", () => {
    categories.scrollBy({ left: -150, behavior: "smooth" });
  });

  rgth_arrow.addEventListener("click", () => {
    categories.scrollBy({ left: 150, behavior: "smooth" });
  });
  const categories = createHTMLel("div", "categories");
  let res = await fetch("/get_categories");
  let data = await res.json();
  data.forEach((category) => {
    let opstion;
    if (type !== "a") {
      opstion = createHTMLel(type, "category", category.name);
    } else {
      opstion = createHTMLel(type, "category", category.name, {
        key: "href",
        value: `/category/${category.name}`,
      });
    }
    categories.append(opstion);
  });
  categoriesSlider.append(left_arrow, categories, rgth_arrow);
};

let offset = 0;

let nomorPosts = false;

export function setupSPA() {
  root.addEventListener("click", (e) => {
    const link = e.target.closest("a");
    if (link && link.origin === location.origin) {
      e.preventDefault();
      if (link.textContent === "Sign In" || link.textContent === "Sign Up") {
        let title = document.head.querySelector("title");
        let logstyle = document.querySelector(".log");
        document.head.removeChild(title);
        document.head.removeChild(logstyle);
        root.innerHTML = "";
      }
      offset = 0;
      nomorPosts = false;
      const postsContainer = document.querySelector(".postscontainer");
      if (postsContainer) postsContainer.innerHTML = "";

      navigateTo(link.pathname);
    }
  });

  window.addEventListener("popstate", () => handleRoute(location.pathname));

  handleRoute(location.pathname);
}

const showPosts = async (path) => {
  if (nomorPosts) return;
  let loading = false;
  if (path == "/") path = "/posts";
  else path = "/api" + path;
  try {
    if (loading) return;
    loading = true;
    let res = await fetch(`${path}?offset=${offset}`);

    if (!res.ok) {
      navigateTo("/signin");
      return;
    }

    let data = await res.json();

    let postsContainer = document.querySelector(".postscontainer");

    creatPosts(postsContainer, data, "append");
    offset += 20;
  } catch (err) {
    const errorEl = createHTMLel("div", "errorEl", "there is no more posts");
    document?.querySelector(".postscontainer")?.appendChild(errorEl);
    nomorPosts = true;
  } finally {
    loading = false;
  }
};

const creatPosts = (container, data, position) => {
  data.forEach((postData) => {
    const postcontainer = createHTMLel("div", "postContainer", "", {
      key: "post-id",
      value: postData.id,
    });
    const postHeader = createHTMLel("a", "link poster", postData.poster, {
      key: "href",
      value: `/${postData.poster}`,
    });
    const creationDate = createHTMLel(
      "span",
      "creationDate",
      formatDate(postData.createdAt)
    );
    const title = createHTMLel("h3", "Posttitle", postData.title);
    const content = createHTMLel("p", "Postcontent", postData.content);
    const like_dislike_containerP = createHTMLel("div", "likeAndDislikeP");
    handleReaction(like_dislike_containerP, "post", postData);

    const postcategories = createcategories(postData.categories);

    postcontainer.append(
      postHeader,
      creationDate,
      title,
      content,
      postcategories,
      like_dislike_containerP
    );
    if (position === "append") container.append(postcontainer);
    if (position === "prepend") container.prepend(postcontainer);
  });
};

function createcategories(categories) {
  const postcategories = createHTMLel("div", "postcategories");
  categories.forEach((category) => {
    const categoryLink = createHTMLel("a", "category", category, {
      key: "href",
      value: `/category/${category}`,
    });
    postcategories.append(categoryLink);
  });

  postcategories.addEventListener("wheel", (e) => {
    e.preventDefault();

    postcategories.scrollBy({
      top: e.deltaY * 0.3,
      behavior: "smooth",
    });
  });

  return postcategories;
}

function handleReaction(container, target, post, userReaction) {
  const likeBtn = createHTMLel("button", "like" + target, "👍");
  const likeSpan = createHTMLel("span", "likesSpan", post.reactions.likes);
  if (post.reactions.action == "like") likeBtn.classList.add("liked");
  const dislikeBtn = createHTMLel("button", "dislike" + target, "👎");
  const dislikeSpan = createHTMLel(
    "span",
    "dislikesSpan",
    post.reactions.dislikes
  );

  if (post.reactions.action == "dislike") dislikeBtn.classList.add("disliked");
  container.append(likeBtn, likeSpan, dislikeBtn, dislikeSpan);

  if (userReaction === "like") {
    likeBtn.classList.add("liked");
  } else if (userReaction === "dislike") {
    dislikeBtn.classList.add("disliked");
  }

  likeBtn.addEventListener("click", () =>
    handleReactionClick(
      "like",
      post.id,
      likeBtn,
      dislikeBtn,
      likeSpan,
      dislikeSpan,
      target
    )
  );

  dislikeBtn.addEventListener("click", () =>
    handleReactionClick(
      "dislike",
      post.id,
      likeBtn,
      dislikeBtn,
      likeSpan,
      dislikeSpan,
      target
    )
  );
}

async function handleReactionClick(
  type,
  id,
  likeBtn,
  dislikeBtn,
  likeSpan,
  dislikeSpan,
  target
) {
  try {
    let response = await fetch(
      `/reactions?target=${target}&id=${id}&action=${type}`,
      {
        method: "POST",
      }
    );

    if (!response.ok) {
      alert("try to react another time");
      return;
    }

    let data = await response.json();

    likeSpan.textContent = data.likes;
    dislikeSpan.textContent = data.dislikes;
    if (data.action === "like") {
      likeBtn.classList.toggle("liked");
      dislikeBtn.classList.remove("disliked");
    } else if (data.action === "dislike") {
      likeBtn.classList.remove("liked");
      dislikeBtn.classList.toggle("disliked");
    }
  } catch (error) {
    console.error("Error handling reaction:", error);
  }
}

const oneday = 60 * 60 * 24;
const onehour = 60 * 60;
const oneminut = 60;

export function formatDate(time) {
  if (!time) time = Date.now() / 1000 - 1;
  let timeText;
  const date = Date.now() / 1000;
  const elapsed = date - time;
  let days = Math.floor(elapsed / oneday);
  let hours = Math.floor((elapsed % oneday) / onehour);
  let minets = Math.floor((elapsed % onehour) / oneminut);
  let seconds = Math.floor(elapsed % oneminut);
  if (days > 0) {
    timeText = `${days}d`;
  } else if (hours > 0) {
    timeText = `${hours}h`;
  } else if (minets > 0) {
    timeText = `${minets}min`;
  } else {
    timeText = `${seconds}s`;
  }
  return timeText;
}

export async function setupPage() {
  if (!document.querySelector(".header")) {
    await navbar();
    const style = createHTMLel("link", "", "", {
      key: "href",
      value: "/frontend/style/post.css",
    });
    searchBar();

    style.rel = "stylesheet";
    const title = createHTMLel("title", "", "Forum");
    document.head.append(style, title);
    addPostPopUp();
    const postsContainer = createHTMLel("div", "postscontainer");
    main.prepend(sidebarLeft, postsContainer, sidebarRight);
    root.appendChild(main);

    await getuser(sidebarLeft);
    upgradeconnection();
    socketEvents();
  }
  showPosts(location.pathname);
}

export const trackscroll = () => {
  document.addEventListener("scrollend", () => {
    if (window.innerHeight + window.scrollY >= root.offsetHeight - 200) {
      showPosts(location.pathname);
    }
  });
};

export const throttle = (func, wait) => {
  let lastCall = 0;
  return function (...args) {
    let now = Date.now();
    if (now - lastCall >= wait) {
      func.apply(this, args);
      lastCall = now;
    }
  };
};
