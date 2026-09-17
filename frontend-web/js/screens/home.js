// Home screen: entry point with three big action cards.
import { el } from "../ui.js";
import { push } from "../nav.js";
import { bookingSearch } from "./booking-search.js";
import { manageBooking } from "./manage-booking.js";
import { adminLogin } from "./admin-login.js";

function actionCard(icon, title, sub, onClick) {
  return el(
    "div",
    { class: "card tappable", onClick },
    el(
      "div",
      { class: "row" },
      el("div", { style: "font-size:28px" }, icon),
      el("div", { class: "grow" }, el("strong", {}, title), el("div", { class: "muted" }, sub)),
      el("div", { style: "color:#999;font-size:20px" }, "›")
    )
  );
}

export function home() {
  const render = () =>
    el(
      "div",
      {},
      el("h2", { class: "title-lg" }, "Mau ke mana?"),
      actionCard("\u{1F68C}", "Pesan tiket", "Pilih rute, jadwal, dan kursi", () => push(bookingSearch())),
      actionCard("\u{1F5D3}️", "Kelola booking", "Cari & reschedule tiket yang ada", () => push(manageBooking())),
      actionCard("\u{1F510}", "Admin", "Dashboard, jadwal, verifikasi", () => push(adminLogin()))
    );
  return { title: "Travel CRM", render };
}
