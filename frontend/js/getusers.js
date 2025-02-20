import { createHTMLel } from "./helpers.js";
import { socket } from "./index.js";

export const getuser = async (sidebar) => {
  const usersContainer = createHTMLel("div", "usersContainer");
  let res = await fetch("/getusers");
  let nickNames = await res.json();
  nickNames.forEach((nickName) => {
    const user = createHTMLel("div", "user");
    const nickname = createHTMLel("span", "nickname", nickName.nickname);
    user.append(nickname);
    const userpopup = addeventToUser(user, nickName.nickname);
    usersContainer.append(user, userpopup);
    sidebar.append(userpopup);
  });
  sidebar.append(usersContainer);
};

function addeventToUser(user, nickname) {
  const div = cratepopUpForUser(nickname);
  user.addEventListener("click", () => {
    let chatshone = document.querySelector("div.showen");
    if (chatshone) {
      chatshone.classList.toggle("showen");
    }
    console.log(chatshone === div);

    if (chatshone === div) {
      div.classList.remove("showen");
    } else {
      div.classList.add("showen");
    }
  });

  return div;
}

function cratepopUpForUser(nickname) {
  const div = createHTMLel("div", "chatContainer");
  div.id = nickname;
  if (nickname === "test") {
    div.classList.add("showen");
  }

  const headercontainer = createHTMLel("div", "headercontainer");
  const header = createHTMLel("h3", "chatheader", nickname);
  const closepanelEl = createHTMLel("span", "close", "X");
  const body = createHTMLel("div", "body");
  const form = createHTMLel("form", "messagesForm");
  headercontainer.append(header, closepanelEl);
  closepanelEl.addEventListener("click", () => {
    const chatContainer = document.querySelector(".showen");
    chatContainer.classList.toggle("showen");
  });

  fetchMsg(form, nickname, body);
  div.append(headercontainer, body, form);
  return div;
}

function fetchMsg(form, to, body) {
  const inpt = createHTMLel("input", "sendmessage", "", {
    key: "placeholder",
    value: "send a message",
  });
  const submitmsg = createHTMLel("button", "submitmsg", "=>");
  form.append(inpt, submitmsg);
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const message = inpt.value.trim();
    if (message) {
      if (sendMessage(to, message)) {
        const messageel = createHTMLel("div", "message sender", message);
        body.append(messageel);
        inpt.value = "";
      }
    }
    if (!message) {
      alert("you can not send an empty message");
    }
  });
}

function sendMessage(receiver, content) {
  let data;
  socket.send(JSON.stringify({ to: receiver, content: content }));
  socket.onmessage = (e) => {
    data = JSON.parse(e.data);
    console.log(data.status);
    if (data.status === "failed") return false;
  };
  return true;
}
