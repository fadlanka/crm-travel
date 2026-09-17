// Booking step 3: enter passenger details, confirm, and create the booking.
import { api } from "../api.js";
import { el, fmtDate, toast, infoRow } from "../ui.js";
import { push } from "../nav.js";
import { ticketScreen } from "./ticket.js";

export function passengerForm(schedule, seat) {
  let submitting = false;
  const container = el("div", { class: "narrow" });
  const name = el("input", { placeholder: "Nama lengkap" });
  const phone = el("input", { placeholder: "Nomor telepon", inputmode: "tel" });
  const email = el("input", { placeholder: "Email (opsional)", type: "email" });

  function paint() {
    container.innerHTML = "";
    container.append(
      el(
        "div",
        { class: "card" },
        infoRow("Rute", `${schedule.origin} → ${schedule.destination}`),
        infoRow("Keberangkatan", fmtDate(schedule.departure_at)),
        infoRow("Kursi", String(seat))
      ),
      el("label", { class: "field mt" }, el("span", {}, "Nama lengkap *"), name),
      el("label", { class: "field" }, el("span", {}, "Nomor telepon *"), phone),
      el("label", { class: "field" }, el("span", {}, "Email (opsional)"), email),
      el("button", { class: "btn", disabled: submitting, onClick: submit }, submitting ? "Memproses…" : "Konfirmasi booking")
    );
  }

  async function submit() {
    if (!name.value.trim() || !phone.value.trim()) { toast("Nama dan telepon wajib diisi"); return; }
    submitting = true; paint();
    try {
      const detail = await api.createBooking({
        passenger_name: name.value.trim(),
        passenger_phone: phone.value.trim(),
        passenger_email: email.value.trim() || null,
        schedule_id: schedule.id,
        seat_number: seat,
      });
      push(ticketScreen(detail, "Booking Berhasil"));
    } catch (e) { toast(String(e)); submitting = false; paint(); }
  }

  return { title: "Data Penumpang", render: () => { paint(); return container; } };
}
