// Manage booking: look up a booking by code, view it, and start a reschedule.
import { api } from "../api.js";
import { el, toast, ticketCard } from "../ui.js";
import { push } from "../nav.js";
import { reschedule } from "./reschedule.js";

export function manageBooking() {
  let detail = null, loading = false;
  const container = el("div", { class: "narrow" });
  const codeInput = el("input", { placeholder: "mis. BK-DEMO01 atau TK-DEMO0001" });

  function paint() {
    container.innerHTML = "";
    container.append(
      el("label", { class: "field" }, el("span", {}, "Kode booking atau tiket"), codeInput),
      el("button", { class: "btn", disabled: loading, onClick: find }, loading ? "Mencari…" : "Cari booking")
    );
    if (detail) {
      container.append(
        el("div", { class: "mt" }, ticketCard(detail)),
        el("button", { class: "btn mt", onClick: () => push(reschedule(detail)) }, "Reschedule perjalanan ini")
      );
    }
  }

  async function find() {
    const code = codeInput.value.trim();
    if (!code) { toast("Masukkan kode booking/tiket"); return; }
    loading = true; detail = null; paint();
    try { detail = await api.getBooking(code); }
    catch (e) { toast(String(e)); }
    finally { loading = false; paint(); }
  }

  return { title: "Kelola Booking", render: () => { paint(); return container; } };
}
