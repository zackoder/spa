import { createHTMLel, fetchData, navigateTo } from "./helpers.js";
import { originalHTML } from "./index.js";
import Validation from "./validation.js";

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

// Registration function
export const signup = async () => {
  let res = await fetch("/getNickName");
  if (res.ok) navigateTo("/");
  else {
    let styleLink = createHTMLel("link", "log", "", {
      key: "href",
      value: "/frontend/style/register.css",
    });
    styleLink.rel = "stylesheet";
    let title = createHTMLel("title", "", "Sign Up");
    document.head.append(styleLink, title);



    const container = createHTMLel("div", "container");
    const form = createHTMLel("form", "", "", { key: "id", value: "registerForm" });
    form.setAttribute("method", "POST");
    form.setAttribute("autocomplete", "off");

    const form_title = createHTMLel("h3", "title", "Create Account");
    form.appendChild(form_title);

    const nicknameBox = createHTMLel("div", "nickname box");
    const nicknameDiv = createHTMLel("div", "", "", { key: "style", value: "width:100%" });
    const nicknameInput = createHTMLel("input", "", "", { key: "type", value: "text" });
    nicknameInput.setAttribute("name", "nickname");
    nicknameInput.setAttribute("id", "nickname");
    nicknameInput.setAttribute("placeholder", "Nickname");
    const errNickname = createHTMLel("span", "", "tst errorr", { key: "id", value: "errnickname" });
    nicknameDiv.append(nicknameInput, errNickname);
    nicknameBox.appendChild(nicknameDiv);
    form.appendChild(nicknameBox);

    const nameBox = createHTMLel("div", "name box grid");
    ["firstname", "lastname"].forEach((field) => {
      const div = createHTMLel("div");
      const input = createHTMLel("input", "", "", { key: "type", value: "text" });
      input.setAttribute("name", field);
      input.setAttribute("id", field);
      input.setAttribute("placeholder", field.split("n").join("N "));
      const errSpan = createHTMLel("span", "", "tst errorr", { key: "id", value: `err${field}` });
      div.append(input, errSpan);
      nameBox.appendChild(div);
    });
    form.appendChild(nameBox);

    const ageGenderBox = createHTMLel("div", "age-gender box grid");
    const genderDiv = createHTMLel("div");
    const genderSelect = createHTMLel("select", "", "", { key: "name", value: "gender" });
    genderSelect.setAttribute("id", "gender");
    ["", "male", "female"].forEach((opt) => {
      const option = createHTMLel("option", "", opt ? opt.charAt(0).toUpperCase() + opt.slice(1) : "Select Your Gender", { key: "value", value: opt });
      genderSelect.appendChild(option);
    });
    const errGender = createHTMLel("span", "", "tst errorr", { key: "id", value: "errgender" });
    genderDiv.append(genderSelect, errGender);
    ageGenderBox.appendChild(genderDiv);

    const birthDiv = createHTMLel("div");
    const birthInput = createHTMLel("input", "", "", { key: "type", value: "date" });
    birthInput.setAttribute("name", "birthdate");
    birthInput.setAttribute("id", "birthdate");
    birthInput.setAttribute("title", "Date of Birth");
    const errBirth = createHTMLel("span", "", "tst errorr", { key: "id", value: "errbirthdate" });
    birthDiv.append(birthInput, errBirth);
    ageGenderBox.appendChild(birthDiv);

    form.appendChild(ageGenderBox);

    const emailBox = createHTMLel("div", "email box");
    const emailDiv = createHTMLel("div", "", "", { key: "style", value: "width:100%" });
    const emailInput = createHTMLel("input", "", "", { key: "type", value: "email" });
    emailInput.setAttribute("name", "email");
    emailInput.setAttribute("id", "email");
    emailInput.setAttribute("placeholder", "Email");
    const errEmail = createHTMLel("span", "", "tst errorr", { key: "id", value: "erremail" });
    emailDiv.append(emailInput, errEmail);
    emailBox.appendChild(emailDiv);
    form.appendChild(emailBox);

    const passwordBox = createHTMLel("div", "password box grid");
    ["password", "confirmpassword"].forEach((field) => {
      const div = createHTMLel("div", "", "", { key: "style", value: "display:block;" });
      const input = createHTMLel("input", "", "", { key: "type", value: "password" });
      input.setAttribute("name", field);
      input.setAttribute("id", field);
      input.setAttribute("placeholder", field === "password" ? "Password" : "Confirm Password");
      const errSpan = createHTMLel("span", "", "tst errorr", { key: "id", value: `err${field}` });
      div.append(input, errSpan);
      passwordBox.appendChild(div);
    });
    form.appendChild(passwordBox);

    const submitBox = createHTMLel("div", "box");
    const submitBtn = createHTMLel("input", "", "", { key: "type", value: "submit" });
    submitBtn.setAttribute("value", "Register");
    submitBtn.setAttribute("id", "submit");
    submitBox.appendChild(submitBtn);
    form.appendChild(submitBox);

    const loginText = createHTMLel("h4", "", "Already have an account? ");
    const loginLink = createHTMLel("a", "", "Login", { key: "href", value: "/login" });
    loginText.appendChild(loginLink);
    form.appendChild(loginText);

    container.appendChild(form);
    root.appendChild(container);


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
        confirmPassword: cpasswordinpt.value.trim(),
      };

      // Validate nickname

      if (!Validation.validateNickname(data.nickName.value.trim())) { return; }
      if (!Validation.validateFirstname(data.firstName.value.trim())) { return; }
      if (!Validation.validateLastname(data.lastName.value.trim())) { return; }
      if (!Validation.validateGender(data.gender)) { return; }
      if (!Validation.validateAge(data.age.value.trim())) { return; }
      if (!Validation.validateEmail(data.email)) { return; }
      if (!Validation.validatePassword(data.password)) { return; }
      if (!Validation.validateConfirmPassword(data.confirmPassword)) { return; }

      let res = fetchData("/sign-up", data);
      res.then((resp) => {
        if (resp.redirected) {
          document.head.removeChild(form_title);
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
    //   form.append(
    //     h1,
    //     nicknameLbl,
    //     nicknameInpt,
    //     firstNamelbl,
    //     firstNameinpt,
    //     lastNamelbl,
    //     lastNameinpt,
    //     agelbl,
    //     ageinpt,
    //     malelbl,
    //     maleipt,
    //     femalelbl,
    //     femaleipt,
    //     Emaillbl,
    //     Emailinpt,
    //     passwordlbl,
    //     passwordinpt,
    //     cpasswordlbl,
    //     cpasswordinpt,
    //     signin,
    //     submitbtn
    //   );
    //   formcontainer.appendChild(form);
    //   root.appendChild(formcontainer);
  }
};

export const signout = async () => {
  try {
    let res = await fetch("/signout");
    console.log("res: ", res);

    if (res.ok) {
      root.innerHTML = "";
      navigateTo("/signin");
    }
  } catch (err) {
    console.log(err);

    alert("an error acursed while signing out");
  }
};
