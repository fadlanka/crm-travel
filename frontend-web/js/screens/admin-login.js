// Admin login. Verifies credentials by calling a protected endpoint; on success
// they are reused as the Basic Auth header for later admin requests.
import { api, setAdminAuth, clearAdminAuth } from "../api.js";
import { el, toast } from "../ui.js";
import { push } from "../nav.js";
import { adminHome } from "./admin-home.js";

export function adminLogin() {
  let loading = false;
  const container = el("div", { class: "narrow" });
  const user = el("input", { value: "admin" });
  const pass = el("input", { type: "password" });

  function paint() {
    container.innerHTML = "";
    container.append(
      el(
        "div",
        { class: "card" },
        el("div", { class: "center", style: "font-size:32px" }, "\u{1F510}"),
        el("label", { class: "field mt" }, el("span", {}, "Username"), user),
        el("label", { class: "field" }, el("span", {}, "Password"), pass),
        el("button", { class: "btn", disabled: loading, onClick: login }, loading ? "Masuk…" : "Masuk"),
        el("div", { class: "center muted mt" }, "Demo: admin / admin123")
      )
    );
  }

  async function login() {
    loading = true; paint();
    setAdminAuth(user.value.trim(), pass.value);
    try { await api.adminDashboard(); push(adminHome()); }
    catch (e) { clearAdminAuth(); toast("Login gagal: " + e); loading = false; paint(); }
  }

  return { title: "Admin Login", render: () => { paint(); return container; } };
}
