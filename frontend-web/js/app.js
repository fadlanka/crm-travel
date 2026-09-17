// Entry point. Wires up the navigator and shows the home screen.
// Every screen lives in its own file under js/screens/.
import { initNav, reset } from "./nav.js";
import { home } from "./screens/home.js";

initNav();
reset(home());
