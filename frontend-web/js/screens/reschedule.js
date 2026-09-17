// Reschedule: pick a new schedule, pick a free seat, confirm.
import { api } from "../api.js";
import { el, fmtDate, toast, spinner, seatCabin } from "../ui.js";
import { push } from "../nav.js";
import { ticketScreen } from "./ticket.js";

export function reschedule(current) {
  let schedules = [], loadingSched = true, newSchedule = null, seatMap = null, newSeat = null, submitting = false;
  const container = el("div", {});

  function schedList() {
    return el(
      "div",
      {},
      ...schedules.map((s) => {
        const full = s.seats_available <= 0;
        const sel = newSchedule && newSchedule.id === s.id;
        return el(
          "div",
          { class: "card tappable" + (full ? " disabled" : ""), style: sel ? "background:#bbdefb" : "", onClick: full ? null : () => pickSchedule(s) },
          el(
            "div",
            { class: "list-row" },
            el("div", {}, el("div", {}, `${s.origin} → ${s.destination}`), el("div", { class: "muted" }, fmtDate(s.departure_at))),
            el("div", { class: "end" }, full ? el("span", { class: "full" }, "PENUH") : `${s.seats_available} kursi`)
          )
        );
      })
    );
  }

  function paint() {
    container.innerHTML = "";
    container.append(
      el(
        "div",
        { class: "card", style: "background:#e3f2fd" },
        el("strong", {}, "Perjalanan sekarang"),
        el("div", {}, `${current.origin} → ${current.destination}`),
        el("div", { class: "muted" }, `${fmtDate(current.schedule.departure_at)} · kursi ${current.booking.seat_number}`)
      ),
      el("div", { class: "title-md" }, "Pilih jadwal baru"),
      loadingSched ? spinner() : schedList()
    );
    if (newSchedule) {
      container.append(
        el("div", { class: "title-md" }, "Pilih kursi"),
        seatMap ? seatCabin(seatMap, newSeat, (s) => { newSeat = s; paint(); }) : spinner()
      );
    }
    container.append(
      el(
        "div",
        { class: "mt" },
        el(
          "button",
          { class: "btn", disabled: !newSchedule || !newSeat || submitting, onClick: confirm },
          submitting ? "Memproses…" : "Konfirmasi reschedule"
        )
      )
    );
  }

  async function pickSchedule(s) {
    newSchedule = s; seatMap = null; newSeat = null; paint();
    try { seatMap = await api.seats(s.id); } catch (e) { toast(String(e)); }
    paint();
  }

  async function confirm() {
    submitting = true; paint();
    try {
      const updated = await api.reschedule(current.booking.booking_code, {
        new_schedule_id: newSchedule.id,
        new_seat_number: newSeat,
      });
      push(ticketScreen(updated, "Reschedule Berhasil"));
    } catch (e) { toast(String(e)); submitting = false; paint(); }
  }

  api.schedules().then((s) => { schedules = s || []; loadingSched = false; paint(); })
    .catch((e) => { toast(String(e)); loadingSched = false; paint(); });

  return { title: "Reschedule", render: () => { paint(); return container; } };
}
