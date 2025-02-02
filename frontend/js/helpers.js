export const addevents = (target, type, path, email, password) => {
  target.addEventListener(type, (e) =>
    fetchSigninData(e, path, email, password)
  );
};

const fetchSigninData = async (e, path, email, password) => {
  e.preventDefault();

  let resp = await fetch(path, {
    method: "POST",
    headers: {
      "Content-type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      email: email,
      password: password,
    }),
  });
  resp.json().then((stract) => console.log(stract.message));
};

export const createHTMLel = (
  name,
  Class,
  content = "",
  atrebute = { key: "", value: "" }
) => {
  let element = document.createElement(name);
  if (content) element.textContent = content;
  if (Class) element.className = Class;
  if (atrebute.key) element.setAttribute(atrebute.key, atrebute.value);
  return element;
};
