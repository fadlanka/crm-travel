// Booking step 2: choose an available seat on the chosen schedule.
import { api } from "../api.js";
import { el, fmtDate, spinner, seatCabin } from "../ui.js";
import { push } from "../nav.js";
import { passengerForm } from "./passenger-form.js";

export function seatSelection(schedule) {
  let seatMap = null, selected = null, loading = true, error = null;
  const container = el("div", { class: "narrow" });

  function paint() {
    container.innerHTML = "";
    if (loading) { container.append(spinner()); return; }
    if (error) { container.append(el("div", { class: "center muted" }, error)); return; }
    container.append(
      el("div", { class: "title-md" }, `${schedule.origin} → ${schedule.destination}`),
      el("div", { class: "muted mb" }, fmtDate(schedule.departure_at)),
      seatCabin(seatMap, selected, (s) => { selected = s; paint(); }),
      el(
        "div",
        { class: "mt" },
        el(
          "button",
          { class: "btn", disabled: selected == null, onClick: () => push(passengerForm(schedule, selected)) },
          selected == null ? "Pilih kursi dulu" : `Lanjut dengan kursi ${selected}`
        )
      )
    );
  }

  api.seats(schedule.id).then((m) => { seatMap = m; loading = false; paint(); })
    .catch((e) => { error = String(e); loading = false; paint(); });

  return { title: "Pilih Kursi", render: () => { paint(); return container; } };
}
