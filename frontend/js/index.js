import { addevents } from "./helpers.js";

const root = document.querySelector(".root");

document.addEventListener("DOMContentLoaded", () => {
  let path = location.pathname;
  if (path === "/signin") {
    signin();
  } else if (path === "/signup") {
    signup();
  }
});

const signin = () => {
  let styleLink = createHTMLel("link");
  styleLink.rel = "stylesheet";
  styleLink.href = "/frontend/style/log.css";
  document.head.appendChild(styleLink);

  /* singin header */
  let h1 = createHTMLel("h1", "logheader", "Sing In");

  /* label and input for email */
  let formcontainer = createHTMLel("div", "formcontainer");

  let form = createHTMLel("form", "logform", "", {
    key: "method",
    value: "POST",
  });

  let Emaillbl = createHTMLel("label", "lbl", "Enter Your Email or Nickname:", {
    key: "for",
    value: "emailnpt",
  });

  let Emailinpt = createHTMLel("input", "inpt", "", {
    key: "id",
    value: "emailnpt",
  });

  // Emailinpt.id = "emailnpt";

  /* label and input for password */
  let passwordlbl = createHTMLel("label", "lbl", "password :", {
    key: "for",
    value: "passwordnpt",
  });

  let passwordinpt = createHTMLel("input", "inpt", "", {
    key: "name",
    value: "password",
  });
  passwordinpt.id = "passwordnpt";
  passwordinpt.type = "password";

  /* submit btn */
  // console.log(path);

  let submitbtn = createHTMLel("button", "btn", "submit");
  addevents(
    submitbtn,
    "click",
    "/sign-in",
    Emailinpt.value,
    passwordinpt.value
  );

  addevents(form, "submit", "/sign-in", Emailinpt.value, passwordinpt.value);

  form.append(h1, Emaillbl, Emailinpt, passwordlbl, passwordinpt, submitbtn);
  formcontainer.appendChild(form);
  root.appendChild(formcontainer);
};

const signup = () => {
  let styleLink = createHTMLel("link", "", "", {
    key: "href",
    value: "/frontend/style/log.css",
  });
  styleLink.rel = "stylesheet";
  // styleLink.href = "/frontend/style/log.css";
  document.head.appendChild(styleLink);

  /* form container */
  let formcontainer = createHTMLel("div", "formcontainer");

  /* form */
  let form = createHTMLel("form", "logform", "", {
    key: "method",
    value: "POST",
  });

  /* Sing Up header */
  let h1 = createHTMLel("h1", "logheader", "Sing Up");

  /* nickname label and  input */
  let nicknameLbl = createHTMLel("label", "lbl", "NickName: ", {
    key: "for",
    value: "nicknameInpt",
  });
  let nicknameInpt = createHTMLel("input", "inpt", "", {
    key: "id",
    value: "nicknameInpt",
  });

  /* first name label and input */
  let firstNamelbl = createHTMLel("label", "lbl", "First Name:", {
    key: "for",
    value: "firstNameinpt",
  });
  let firstNameinpt = createHTMLel("input", "inpt", "", {
    key: "id",
    value: "firstNameinpt",
  });

  /* last name label and input */
  let lastNamelbl = createHTMLel("label", "lbl", "Last Name:", {
    key: "for",
    value: "lastNameinpt",
  });
  let lastNameinpt = createHTMLel("input", "inpt", "", {
    key: "id",
    value: "lastNameinpt",
  });

  /* age label and input */
  let agelbl = createHTMLel("label", "lbl", "Age: ", {
    key: "for",
    value: "ageinpt",
  });

  let ageinpt = createHTMLel("input", "inpt", "", {
    key: "id",
    value: "ageinpt",
  });
  ageinpt.type = "date";

  form.append(
    h1,
    nicknameLbl,
    nicknameInpt,
    firstNamelbl,
    firstNameinpt,
    lastNamelbl,
    lastNameinpt,
    agelbl,
    ageinpt
  );
  formcontainer.appendChild(form);
  root.appendChild(formcontainer);
};

function createHTMLel(
  name,
  Class,
  content = "",
  atrebute = { key: "", value: "" }
) {
  let element = document.createElement(name);
  if (content) element.textContent = content;
  if (Class) element.className = Class;
  if (atrebute.key) element.setAttribute(atrebute.key, atrebute.value);
  return element;
}

function test(...args) {
  console.log(args);
}
