// Screen-stack navigator (like Flutter's Navigator).
//
// The app is a stack of screens. push() shows a new screen on top, pop() (the
// ← button) removes the top one, reset() clears the stack and starts fresh.
// A "screen" is any object { title, render() } where render() returns a DOM node.

const stack = [];
let appEl, titleEl, backBtn;

// Wire up the DOM elements and the back button. Call once at startup.
export function initNav() {
  appEl = document.getElementById("app");
  titleEl = document.getElementById("title");
  backBtn = document.getElementById("backBtn");
  backBtn.addEventListener("click", pop);
}

// Render whatever screen is currently on top of the stack.
function renderTop() {
  const view = stack[stack.length - 1];
  appEl.innerHTML = "";
  appEl.append(view.render());
  titleEl.textContent = view.title;
  backBtn.hidden = stack.length <= 1; // hide ← on the first screen
}

export function push(view) {
  stack.push(view);
  renderTop();
}

// Only the ← button pops, so this stays internal to the navigator.
function pop() {
  if (stack.length > 1) {
    stack.pop();
    renderTop();
  }
}

export function reset(view) {
  stack.length = 0;
  stack.push(view);
  renderTop();
}
