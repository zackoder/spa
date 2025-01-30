const navbar = document.querySelector(".nav");
const navLinks = document.querySelectorAll("[data-link]");

// console.log(navLinks);

// navbar.addEventListener("click", (e) => {
//   let target = e.target;

//   if (target.matches("[data-link]")) {
//     target.preventDefault(); /* preventDefault */
//     console.log(target);
//   }
// });

const router = async () => {
  const routes = [
    { path: "/NotFound", isMatch: false },
    { path: "/", isMatch: true },
    { path: "/singin", isMatch: false },
    { path: "/singup", isMatch: false },
    { path: "/profile", isMatch: false },
  ];
  const potentialmatches = routes.map((route) => {
    return {
      route: route,
      isMatch: location.pathname === route.path,
    };
  });
  let match = potentialmatches.find(
    (potentialmatche) => potentialmatche.isMatch
  );
  if (!match) {
    match = {
      route: routes[0],
      isMatch: true,
    };
  }
  console.log(match);
};

document.addEventListener("DOMContentLoaded", () => {
  router();
});

function createHTMLel(
  name,
  atrebute = { key: "", value: "" },
  Class = "",
  content = ""
) {
  let element = document.createElement(name);
  if (content) element.textContent = content;
  if (Class) element.classlist = Class;
  if (atrebute.key) element.setAttribute(atrebute.key, atrebute.value);
  return element;
}

let test = createHTMLel("a", { key: "data-link", value: "" }, "", "Profile");
test.href = "/profile";
let br = createHTMLel("br");
navbar.append(br, test);
