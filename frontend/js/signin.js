import { createHTMLel, fetchData } from "./helpers.js";

const root = document.querySelector(".root");

export const signin = () => {
  let styleLink = createHTMLel("link", "", "", {
    key: "href",
    value: "/frontend/style/log.css",
  });

  styleLink.rel = "stylesheet";

  let title = createHTMLel("title", "", "Sign In");

  document.head.append(styleLink, title);

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

  let submitbtn = createHTMLel("button", "submit", "submit");
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    let email = Emailinpt.value.trim();
    let password = passwordinpt.value.trim();
    const data = {
      email: email,
      password: password,
    };

    let res = fetchData("/sign-in", data);
    res.then((res) => {
      if (res.ok) location.href = "/";
    });
  });

  form.append(h1, Emaillbl, Emailinpt, passwordlbl, passwordinpt, submitbtn);
  formcontainer.appendChild(form);
  root.appendChild(formcontainer);
};

export const signup = () => {
  let styleLink = createHTMLel("link", "", "", {
    key: "href",
    value: "/frontend/style/log.css",
  });
  styleLink.rel = "stylesheet";
  let title = createHTMLel("title", "", "Sign Up");
  document.head.append(styleLink, title);

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

  let malelbl = createHTMLel("label", "lbl", "Male", {
    key: "for",
    value: "male",
  });

  let maleipt = createHTMLel("input", "redioinpt", "", {
    key: "name",
    value: "gender",
  });

  maleipt.type = "radio";
  maleipt.id = "male";

  let femalelbl = createHTMLel("label", "lbl", "Female", {
    key: "for",
    value: "female",
  });
  let femaleipt = createHTMLel("input", "redioinpt", "", {
    key: "name",
    value: "gender",
  });

  femaleipt.type = "radio";
  femaleipt.id = "female";

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

  /* conferm password */
  let cpasswordlbl = createHTMLel("label", "lbl", "conferm password :", {
    key: "for",
    value: "cpasswordnpt",
  });

  let cpasswordinpt = createHTMLel("input", "inpt", "", {
    key: "name",
    value: "password",
  });
  passwordinpt.id = "cpasswordnpt";
  passwordinpt.type = "password";

  let submitbtn = createHTMLel("button", "submit", "submit");

  form.addEventListener("submit", (e) => {
    e.preventDefault();

    let email = Emailinpt.value.trim();
    let password = passwordinpt.value.trim();
    let gender = maleipt.checked ? "male" : femaleipt.checked ? "female" : "";
    const data = {
      nickName: nicknameInpt.value,
      firstName: firstNameinpt.value,
      lastName: lastNameinpt.value,
      gender: gender,
      age: ageinpt.value,
      email: email,
      password: password,
    };

    let res = fetchData("/signup", data);
    res.then((resp) => {
      console.log(resp);
      if (resp.redirected) {
        location.href = "/";
      }
    });
  });

  form.append(
    h1,
    nicknameLbl,
    nicknameInpt,
    firstNamelbl,
    firstNameinpt,
    lastNamelbl,
    lastNameinpt,
    agelbl,
    ageinpt,
    malelbl,
    maleipt,
    femalelbl,
    femaleipt,
    Emaillbl,
    Emailinpt,
    passwordlbl,
    passwordinpt,
    cpasswordlbl,
    cpasswordinpt,
    submitbtn
  );
  formcontainer.appendChild(form);
  root.appendChild(formcontainer);
};

export const signout = async () => {
  try {
    let res = await fetch("/signout");
    console.log(res);

    if (res.ok) {
      location.href = "/signin";
    }
  } catch (err) {
    alert("an error acursed while signing out");
  }
};
