import { root } from "./navbar.js";

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

const main = createHTMLel("main", "main");
const sidebarLeft = createHTMLel("aside", "sidebar left-sidebar");
const sidebarRight = createHTMLel("aside", "sidebar right-sidebar");

// export const addevents = (target, type, path, data) => {
// target.addEventListener(type, (e) => fetchData(e, path, data));
// };

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
  console.log(data);
  const res = await fetch("/addpost", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
};

export const addPostPopUp = () => {
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
  creatcategories(categories, "div");

  let chosenCategoreis = Array.from(
    categories.querySelectorAll(".categories .category")
  );
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

export const creatcategories = (categoriesSlider, type) => {
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
  for (let i = 0; i < 50; i++) {
    const opstion = createHTMLel(type, "category", "test", {
      key: "href",
      value: `/category/test`,
    });
    categories.append(opstion);
  }
  categoriesSlider.append(left_arrow, categories, rgth_arrow);
};

let offset = 0;

export const showPosts = async (path) => {
  let loading = false;
  try {
    if (loading) return;
    loading = true;
    let res = await fetch(`${path}?offset=${offset}`);

    let data = await res.json();
    let postsContainer = document.querySelector(".postscontainer");

    if (!postsContainer) {
      postsContainer = createHTMLel("div", "postscontainer");
      main.append(sidebarLeft, postsContainer, sidebarRight);
      root.appendChild(main);
    }

    creatPosts(postsContainer, data);
    console.log("hello");
    offset += 20;
  } catch {
  } finally {
    loading = false;
  }
};

const creatPosts = (container, data) => {
  data.forEach((postData) => {
    const postcontainer = createHTMLel("div", "postContainer", "", {
      key: "post-id",
      value: postData.id,
    });
    const postHeader = createHTMLel("h2", "poster", postData.poster);
    const creationDate = createHTMLel(
      "span",
      "creationDate",
      formatDate(postData.createdAt)
    );
    const title = createHTMLel("h3", "Posttitle", postData.title);
    const content = createHTMLel("p", "Postcontent", postData.content);
    const like_dislike_containerP = createHTMLel("div", "likeAndDislikeP");
    handleReaction(like_dislike_containerP, "post", postData);
    postcontainer.addEventListener("click", (e) => {
      console.log(e.target);
    });
    postcontainer.append(
      postHeader,
      creationDate,
      title,
      content,
      like_dislike_containerP
    );
    container.append(postcontainer);
  });
};

function handleReaction(container, target, post) {
  const likebtn = createHTMLel("button", "like" + target, "like");
  const likespan = createHTMLel("span", "likesSpan", post.reactions.likes);
  console.log(post.reactions.likes);

  const dislikebtn = createHTMLel("button", "dislike" + target, "dislike");
  const dislikespan = createHTMLel(
    "span",
    "dislikesSpan",
    post.reactions.dislikes
  );
  container.append(likebtn, likespan, dislikebtn, dislikespan);
}

const oneday = 60 * 60 * 24;
const onehour = 60 * 60;
const oneminut = 60;

function formatDate(time) {
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

export const handlescroll = (scroll) => {
  let currentscroll = window.scrollY;

  if (scroll < currentscroll) {
    showPosts("/posts");
    currentscroll = scroll;
  }
};
