import { root } from "./navbar.js";
export const addevents = (target, type, path, email, password) => {
  target.addEventListener(type, (e) =>
    fetchSigninData(e, path, email, password)
  );
};

const fetchSigninData = async (e, path, email, password) => {
  e.preventDefault();

  let resp = await fetch(path, {
    method: "POST",
    headers: {
      "Content-type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      email: email,
      password: password,
    }),
  });
  resp.json().then((stract) => console.log(stract.message));
};

export const createHTMLel = (
  name,
  Class,
  content = "",
  atrebute = { key: "", value: "" }
) => {
  let element = document.createElement(name);
  if (content) element.textContent = content;
  if (Class) element.className = Class;
  if (atrebute.key) element.setAttribute(atrebute.key, atrebute.value);
  return element;
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
    console.log(path);

    if (loading) return;
    loading = true;
    let res = await fetch(`${path}?offset=${offset}`);
    let data = await res.json();
    let postsContainer = document.querySelector(".postscontainer");

    if (!postsContainer) {
      postsContainer = createHTMLel("div", "postscontainer");
      root.appendChild(postsContainer);
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
    console.log(postData);

    const postcontainer = createHTMLel("div", "postContainer");
    const title = createHTMLel("h2", "Posttitle", postData.title);
    const content = createHTMLel("p", "Postcontent", postData.content);
    postcontainer.addEventListener("click", (e) => {});
    postcontainer.append(title, content);
    container.append(postcontainer);
  });
};
