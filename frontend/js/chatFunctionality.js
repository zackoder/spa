import { createmsgcontaine } from "./getusers.js";
import { createHTMLel } from "./helpers.js";
import { user } from "./navbar.js";

export let socket = null;

export const socketEvents = () => {
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

    console.log("sender element", senderchatbox);

    let newMessage;

    if (senderchatbox && !senderchatbox.classList.contains("showen")) {
      console.log("it is hiden");

      if (senderchatbox.children[1].children.length !== 0) {
        console.log("it has cheldren");
        newMessage = createmsgcontaine(data, data.from);
        newMessage.classList.add("get");
        senderchatbox.children[1].prepend(newMessage);
        scrolldown(parent, newMessage);
      }
      document.querySelectorAll(`.user`).forEach((user) => {

        if (user.children[0].textContent === data.from) {
          const notificationEl = createHTMLel("span", "notification");
          console.log("hello");
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

  socket.addEventListener("open", () => {
    console.log("connected..");
  });
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
    users.forEach((user) => {
      if (user.children[0].textContent === data.nickname)
        user.children[1].textContent = "online";
    });
  } else if (data.user === "offline") {
    users.forEach((user) => {
      if (user.children[0].textContent === data.nickname)
        user.children[1].textContent = "offline";
    });
  }
}
