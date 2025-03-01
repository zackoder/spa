import { createmsgcontaine } from "./getusers.js";
import { createHTMLel } from "./helpers.js";

export let socket = null;

export const socketEvents = () => {
  socket.onopen = (e) => {
    console.log("the client is connected to the server");
  };

  socket.onmessage = (e) => {
    const data = JSON.parse(e.data);
    console.log(data);

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
      }
      const notifacation = document
        .querySelectorAll(`.user`)
        .forEach((user) => {
          if (user.textContent === data.from) {
            const notificationEl = createHTMLel("span", "notification");
            user.append(notificationEl);
          }
        });
    } else {
      newMessage = createmsgcontaine(data, data.from);

      newMessage.classList.add("get");

      senderchatbox.children[1].prepend(newMessage);
    }

    setTimeout(() => {
      parent.scrollTop = parent.scrollHeight + newMessage.offsetHeight;
    }, 0);
  };

  socket.addEventListener("open", () => {
    console.log("connected..");
  });
};

export function upgradeconnection() {
  if (socket !== null) return;
  socket = new WebSocket("ws://localhost:8080/ws");
}
