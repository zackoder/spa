const navbar = document.querySelector(".nav");
const navLinks = document.querySelectorAll("[data-link]");

console.log(navLinks);

navbar.addEventListener("click", (e) => {
  let target = e.target;

  if (target.matches("[data-link]")) {
    target.preventDefault(); /* preventDefault */
    console.log(target);
  }
});

const router = async () => {
  const routes = [
    { path: "/", view: "/" },
    { path: "/singin", view: "/singin" },
    { path: "/singup", view: "/singup" },
  ];
  const potentialmatches = routes.map((route) => {
    return {
      route: route,
      isMatch: location.pathname === route.path,
    };
  });
};
