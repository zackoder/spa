import { createHTMLel, formatDate, throttle } from "./helpers.js";
import { socket } from "./chatFunctionality.js";
import { user } from "./navbar.js";

export const getuser = async (sidebar) => {
  const usersContainer = createHTMLel("div", "usersContainer");
  let res = await fetch("/getusers");
  let nickNames = await res.json();
  nickNames.forEach((nickName) => {
    const user = createHTMLel("div", "user");
    const nickname = createHTMLel("span", "nickname", nickName.nickname);
    const stat = createHTMLel("span", "stat", "offline");
    user.append(nickname, stat);
    const userpopup = addeventToUser(user, nickName.nickname);
    usersContainer.append(user);
    sidebar.append(userpopup);
  });
  sidebar.append(usersContainer);
};

let messagesOffset = 0;

function addeventToUser(user, nickname) {
  const div = cratepopUpForUser(nickname);
  user.addEventListener("click", async () => {
    let chatshone = document.querySelector("div.showen");
    let notifacation = user.querySelector(".notification");
    if (notifacation) user.removeChild(notifacation);
    if (chatshone && chatshone !== div) {
      chatshone.classList.toggle("showen");
    }
    messagesOffset = 0;
    const body = div.children[1];
    div.classList.toggle("showen");
    const data = await getmessages(nickname);
    if (body.children.length === 0) {
      if (data) {
        data.forEach((msg) => {
          creatmessage(msg, div.querySelector(".body"), nickname, "append");
        });
      }
    }

    const throttledScrollHandler = throttle(
      () => handelmessagesscroll(nickname),
      1000
    );
    body.addEventListener("scroll", throttledScrollHandler);
  });
  return div;
}

async function getmessages(nickname) {
  const res = await fetch(
    `/api/messages?offset=${messagesOffset}&to=${nickname}`
  );
  const data = await res.json();

  if (data === null) return;
  else {
    messagesOffset += data.length;
    return data;
  }
}

export function creatmessage(msg, parent, nickName, possition) {
  const div = createmsgcontaine(msg);

  if (msg.to === user) {
    div.children[0].textContent = nickName;
    div.classList.add("get");
  } else if (msg.to === nickName) {
    div.children[0].textContent = user;
    div.classList.add("sent");
  }

  parent.append(div);

  setTimeout(() => {
    parent.scrollTop = parent.scrollTop;
  }, 0);
}

export const createmsgcontaine = (msg, name) => {
  const div = createHTMLel("div", "message");
  const sender = createHTMLel("h4", "sender", name);
  const content = createHTMLel("span", "content", msg.content);
  const date = formatDate(msg.creationDate);
  div.append(sender, content, date);
  return div;
};

function cratepopUpForUser(nickname) {
  const div = createHTMLel("div", "chatContainer");
  div.id = nickname;

  const headercontainer = createHTMLel("div", "headercontainer");
  const header = createHTMLel("h3", "chatheader", nickname);
  const closepanelEl = createHTMLel("span", "close", "X");
  const body = createHTMLel("div", "body");
  const errorEl = createHTMLel("p", "errorp");
  const form = createHTMLel("form", "messagesForm");
  headercontainer.append(header, closepanelEl);
  closepanelEl.addEventListener("click", () => {
    const chatContainer = document.querySelector(".showen");
    chatContainer.classList.toggle("showen");
  });
  const inpt = createHTMLel("input", "sendmessage", "", {
    key: "placeholder",
    value: "send a message",
  });
  const submitmsg = createHTMLel("button", "submitmsg", "=>");
  form.append(inpt, submitmsg);

  fetchMsg(form, nickname, errorEl, body);
  div.append(headercontainer, body, errorEl, form);
  return div;
}

async function handelmessagesscroll(nickname) {
  const element = document.querySelector(`#${nickname} .body`);
  const elrect = element.getBoundingClientRect();

  const data = await getmessages(nickname);
  if (
    Math.abs(element.scrollTop) >=
    element.scrollHeight - (elrect.height + 100)
  ) {
    if (data.length > 0) {
      data.forEach((msg) => {
        creatmessage(msg, element, nickname, "prepend");
      });
    }
  }
}

function fetchMsg(form, to, errEl, body) {
  let inpt = form.children[0];
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const message = inpt.value.trim();

    if (message) {
      let res = await sendMessage(to, message);
      if (res.status === "failed") {
        errEl.textContent = res.message;
        return;
      }
      if (res.status === "successe") {
        if (errEl.textContent) errEl.textContent = "";
        const messageEl = createmsgcontaine(
          {
            content: message,
            creationDate: 0,
          },
          user
        );
        messageEl.classList.add("sent");
        body.prepend(messageEl);
        let users = document.querySelector(".usersContainer");
        document.querySelectorAll(".user").forEach((user) => {
          if (user.children[0].textContent === to) {
            users.prepend(user);
          }
        });
        messagesOffset++;
        inpt.value = "";
      }
    }

    if (!message) {
      alert("you can not send an empty message");
    }
  });
}

async function sendMessage(receiver, content) {
  return new Promise((res) => {
    let data;
    socket.send(JSON.stringify({ to: receiver, content: content }));
    socket.onmessage = (e) => {
      data = JSON.parse(e.data);
      if (data.from) {
        const senderchatbox = document.querySelector("#" + data.from);

        if (
          senderchatbox !== null &&
          senderchatbox.classList.contains(".showen")
        ) {

          const newMessage = createmsgcontaine(data, data.from);
          newMessage.classList.add("get");
          senderchatbox.children[1].prepend(newMessage);
          setTimeout(() => {
            parent.scrollTop =
              parent.scrollHeight + newMessage.offsetHeight + 10;
          }, 10);
          return;
        } else {

          if (senderchatbox.children[1].children.length !== 0) {
            const newMessage = createmsgcontaine(data, data.from);
            newMessage.classList.add("get");
            senderchatbox.children[1].prepend(newMessage);
          }
        }

        document.querySelectorAll(`.user`).forEach((user) => {
          if (user.children[0].textContent === data.from) {
            const notificationEl = createHTMLel("span", "notification");
            user.append(notificationEl);
          }
        });
      }

      res(data);
    };
  });
}
