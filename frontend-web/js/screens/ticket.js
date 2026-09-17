// Ticket screen: shows the digital ticket (with QR) after booking/reschedule.
import { el, ticketCard } from "../ui.js";
import { reset } from "../nav.js";
import { home } from "./home.js";

export function ticketScreen(detail, title = "Tiket") {
  const render = () =>
    el(
      "div",
      { class: "narrow" },
      ticketCard(detail),
      el("button", { class: "btn secondary mt", onClick: () => reset(home()) }, "Kembali ke beranda")
    );
  return { title, render };
}
