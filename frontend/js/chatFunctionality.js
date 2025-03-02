import {
  addeventToUser,
  createmsgcontaine,
  createUsrContainer,
} from "./getusers.js";
import { createHTMLel } from "./helpers.js";
import { user } from "./navbar.js";

export let socket = null;

export const socketEvents = () => {
  // socket.send(JSON.stringify({ to: receiver, content: content }));

  socket.onopen = (e) => {
    console.log("the client is connected to the server");
  };

  socket.onmessage = (e) => {
    const data = JSON.parse(e.data);
    if (data.user) {
      handleconnection(data);
      return;
    }

    const senderchatbox = document.querySelector("#" + data.from);

    let newMessage;

    if (senderchatbox && !senderchatbox.classList.contains("showen")) {
      if (senderchatbox.children[1].children.length !== 0) {
        newMessage = createmsgcontaine(data, data.from);
        newMessage.classList.add("get");
        senderchatbox.children[1].prepend(newMessage);
        scrolldown(parent, newMessage);
      }
      const users = document.querySelector(".usersContainer");
      document.querySelectorAll(`.user`).forEach((user) => {
        if (user.children[0].textContent === data.from) {
          const notificationEl = createHTMLel("span", "notification");
          users.prepend(user);
          user.append(notificationEl);
        }
      });
    } else {
      newMessage = createmsgcontaine(data, data.from);

      newMessage.classList.add("get");

      senderchatbox.children[1].prepend(newMessage);
      scrolldown(parent, newMessage);
    }
  };
};

function scrolldown(parent, newMessage) {
  setTimeout(() => {
    parent.scrollTop = parent.scrollHeight + newMessage.offsetHeight;
  }, 0);
}

export function upgradeconnection() {
  if (socket !== null) return;
  socket = new WebSocket("ws://localhost:8080/ws");
}

function handleconnection(data) {
  let users = document.querySelectorAll(".user");
  if (data.user === "online") {
    const getUser = document.querySelector(`#${data.nickname}`);
    if (data.nickname !== user && !getUser) {
      const newUser = createUsrContainer(data, "online");
      const userpopu = addeventToUser(newUser, data.nickname);
      document.querySelector(".left-sidebar").append(userpopu);
      document.querySelector(".usersContainer").append(newUser);
    }
    changeUserstat(users, "online", data.nickname);
  } else if (data.user === "offline") {
    changeUserstat(users, "", data.nickname);
  }
}

function changeUserstat(users, stat, nickname) {
  users.forEach((user) => {
    if (user.children[0].textContent === nickname)
      user.children[1].textContent = stat;
  });
}
