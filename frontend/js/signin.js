import { socket } from "./chatFunctionality.js";
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
      value: "/frontend/style/login.css",
    });

    styleLink.rel = "stylesheet";

    let docTitle = createHTMLel("title", "", "Sign In");

    document.head.append(styleLink, docTitle);

    /* singin header */
    const container = createHTMLel("div", "container");
    const form = createHTMLel("form", "", "", { key: "id", value: "loginForm" });
    form.setAttribute("method", "POST");
    form.setAttribute("autocomplete", "off");

    const title = createHTMLel("h3", "title", "Login");
    form.appendChild(title);

    const usernameDiv = createHTMLel("div");
    const usernameLabel = createHTMLel("label", "", "Nickname / Email", { key: "for", value: "username" });
    const usernameInput = createHTMLel("input", "", "", { key: "type", value: "text" });
    usernameInput.setAttribute("name", "username");
    usernameInput.setAttribute("id", "username");
    const errUsername = createHTMLel("span", "", "testerr", { key: "id", value: "errusername" });
    usernameDiv.append(usernameLabel, usernameInput, errUsername);
    form.appendChild(usernameDiv);

    const passwordDiv = createHTMLel("div", "password");
    const passwordLabel = createHTMLel("label", "", "Password", { key: "for", value: "password" });
    const passwordInnerDiv = createHTMLel("div");
    const passwordInput = createHTMLel("input", "password", "", { key: "type", value: "password" });
    passwordInput.setAttribute("name", "password");
    passwordInput.setAttribute("id", "password");
    const showPasswordBtn = createHTMLel("button", "showpassword", "Show", { key: "id", value: "showpassword" });
    //change input type
    changeInputType(passwordInput, showPasswordBtn);
    passwordInnerDiv.append(passwordInput, showPasswordBtn);
    const errPassword = createHTMLel("span", "", "testerr", { key: "id", value: "errpassword" });
    passwordDiv.append(passwordLabel, passwordInnerDiv, errPassword);
    form.appendChild(passwordDiv);

    // const rememberDiv = createHTMLel("div");
    // const rememberInput = createHTMLel("input", "", "", { key: "type", value: "checkbox" });
    // rememberInput.setAttribute("name", "rememberme");
    // rememberInput.setAttribute("id", "rememberme");
    // const rememberLabel = createHTMLel("label", "rememberme", "Remember me", { key: "for", value: "rememberme" });
    // rememberDiv.append(rememberInput, rememberLabel);
    // form.appendChild(rememberDiv);

    const submitDiv = createHTMLel("div");
    const submitBtn = createHTMLel("input", "", "", { key: "type", value: "submit" });
    submitBtn.setAttribute("value", "Login");
    const registerText = createHTMLel("h3", "", "Don't have an account? ");
    const registerLink = createHTMLel("a", "", "Register", { key: "href", value: "/signup" });
    registerText.appendChild(registerLink);
    submitDiv.append(submitBtn, registerText);
    form.appendChild(submitDiv);

    container.appendChild(form);

    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      let email = usernameInput.value.trim();
      let password = passwordInput.value.trim();
      const data = {
        email: email,
        password: password,
      };

      let res = fetchData("/sign-in", data);
      res.then((res) => {
        if (res.ok) {
          document.head.removeChild(styleLink);
          document.head.removeChild(docTitle);
          root.innerHTML = "";
          navigateTo("/");
        }
      });
    });

    root.appendChild(container);
  }
};

function changeInputType(passwordInput, showPasswordBtn) {
  let show = false

  showPasswordBtn.addEventListener("click", (e) => {
    e.preventDefault()
    if (!show) {
      passwordInput.type = "text";
      showPasswordBtn.textContent = "Hide";
    } else {
      passwordInput.type = "password";
      showPasswordBtn.textContent = "Show";
    }
    show = !show
  });
}

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
    ["Firstname", "Lastname"].forEach((field) => {
      const div = createHTMLel("div");
      const input = createHTMLel("input", field, "", { key: "type", value: "text" });
      input.setAttribute("name", field);
      input.setAttribute("id", field);
      input.setAttribute("placeholder", field.split("n").join(" N"));
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
      const div = createHTMLel("div", field, "", { key: "style", value: "display:block;" });
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
    const loginLink = createHTMLel("a", "", "Login", { key: "href", value: "/signin" });
    loginText.appendChild(loginLink);
    form.appendChild(loginText);

    container.appendChild(form);
    root.appendChild(container);


    form.addEventListener("submit", (e) => {
      e.preventDefault();

      const firstNameinpt = document.querySelector('.Firstname')
      const lastNameinpt = document.querySelector('.Lastname')
      const password = document.querySelector('#password')
      console.log(password);

      const confirmpassword = document.querySelector('#confirmpassword')
      const data = {
        nickName: nicknameInput.value.trim(),
        firstName: firstNameinpt.value.trim(),
        lastName: lastNameinpt.value.trim(),
        gender: genderSelect.value.trim(),
        age: birthInput.value.trim(),
        email: emailInput.value.trim(),
        password: password.value.trim(),
        confirmPassword: confirmpassword.value.trim(),
      };

      // Validate nickname
      if (!Validation.validateNickname(data.nickName)) { return; }
      if (!Validation.validateFirstname(data.firstName)) { return; }
      if (!Validation.validateLastname(data.lastName)) { return; }
      if (!Validation.validateGender(data.gender)) { return; }
      if (!Validation.validateAge(data.age)) { return; }
      if (!Validation.validateEmail(data.email)) { return; }
      if (!Validation.validatePassword(data.password)) { return; }
      if (!Validation.validateConfirmPassword(data.confirmPassword)) { return; }

      let res = fetchData("/sign-up", data);
      res.then((resp) => {
        if (resp.ok) {
          document.head.removeChild(title);
          document.head.removeChild(styleLink);
          root.innerHTML = "";
          navigateTo("/signin");
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
      socket.close();

      navigateTo("/signin");
    }
  } catch (err) {
    console.log(err);

    alert("an error acursed while signing out");
  }
};
