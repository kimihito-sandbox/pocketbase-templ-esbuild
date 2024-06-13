import * as htmx from "htmx.org";
import Cookies from "js-cookie";

window.htmx = htmx;

document.addEventListener("htmx:configRequest", function (event) {
  event.detail.headers["X-CSRF-Token"] = Cookies.get("_csrf");
});


console.log("fooo");
