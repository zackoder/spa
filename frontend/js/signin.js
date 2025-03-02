import { createHTMLel, fetchData, navigateTo } from "./helpers.js";
import { originalHTML } from "./index.js";

const root = document.querySelector(".root");

export const signin = async () => {
  let res = await fetch("/getNickName");
  if (res.ok) navigateTo("/");
  else {
    // document.documentElement.innerHTML = originalHTML;
    let styleLink = createHTMLel("link", "log", "", {
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

    let Emaillbl = createHTMLel(
      "label",
      "lbl",
      "Enter Your Email or Nickname:",
      {
        key: "for",
        value: "emailnpt",
      }
    );

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
        if (res.ok) {
          document.head.removeChild(styleLink);
          document.head.removeChild(title);
          root.innerHTML = "";
          navigateTo("/");
        }
      });
    });

    const signup = createHTMLel("p", "signuplikn", "click the link to ");
    const signupLink = createHTMLel("a", "link", "Sign Up", {
      key: "href",
      value: "/signup",
    });
    signup.appendChild(signupLink);
    form.append(
      h1,
      Emaillbl,
      Emailinpt,
      passwordlbl,
      passwordinpt,
      signup,
      submitbtn
    );
    formcontainer.appendChild(form);
    root.appendChild(formcontainer);
  }
};

export const signup = async () => {
  let res = await fetch("/getNickName");
  if (res.ok) navigateTo("/");
  else {
    let styleLink = createHTMLel("link", "log", "", {
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

    let Emaillbl = createHTMLel(
      "label",
      "lbl",
      "Enter Your Email or Nickname:",
      {
        key: "for",
        value: "emailnpt",
      }
    );

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

      let res = fetchData("/sign-up", data);
      res.then((resp) => {
        if (resp.redirected) {
          document.head.removeChild(title);
          document.head.removeChild(styleLink);
          navigateTo("/");
        }
      });
    });
    const signin = createHTMLel("p", "signuplikn", "click the link to ");
    const signinLink = createHTMLel("a", "link", "Sign In", {
      key: "href",
      value: "/signin",
    });
    signin.appendChild(signinLink);
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
      signin,
      submitbtn
    );
    formcontainer.appendChild(form);
    root.appendChild(formcontainer);
  }
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
