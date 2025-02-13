import { createHTMLel } from "./helpers.js";
import { root } from "./navbar.js";

export const notFound = () => {
  root.innerHTML = "";
  const div = createHTMLel(
    "div",
    "err",
    "the content you are looking for doesn't exists at the moment"
  );
  root.append(div);
};
