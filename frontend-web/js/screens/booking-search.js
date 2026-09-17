// Booking step 1: choose a route (origin/destination) and pick a schedule.
import { api } from "../api.js";
import { el, fmtDate, toast, spinner } from "../ui.js";
import { push } from "../nav.js";
import { seatSelection } from "./seat-selection.js";

// A labelled <select> dropdown.
function selectField(label, value, options, onChange) {
  const sel = el("select", { onChange: (e) => onChange(e.target.value) });
  sel.append(el("option", { value: "" }, `— ${label} —`));
  for (const o of options) {
    const opt = el("option", { value: o }, o);
    if (o === value) opt.selected = true;
    sel.append(opt);
  }
  return el("label", { class: "field grow" }, el("span", {}, label), sel);
}

export function bookingSearch() {
  let routes = [], schedules = [], origin = "", destination = "";
  let loadingRoutes = true, loadingSchedules = false;
  const container = el("div", {});

  function scheduleList() {
    if (!schedules.length)
      return el("div", { class: "center muted", style: "padding:24px" }, "Belum ada jadwal. Coba cari.");
    return el(
      "div",
      {},
      ...schedules.map((s) => {
        const full = s.seats_available <= 0;
        return el(
          "div",
          { class: "card tappable" + (full ? " disabled" : ""), onClick: full ? null : () => push(seatSelection(s)) },
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
    if (loadingRoutes) { container.append(spinner()); return; }
    const origins = [...new Set(routes.map((r) => r.origin))].sort();
    const dests = [...new Set(routes.filter((r) => !origin || r.origin === origin).map((r) => r.destination))].sort();
    container.append(
      el(
        "div",
        { class: "row" },
        selectField("Asal", origin, origins, (v) => { origin = v; destination = ""; paint(); }),
        selectField("Tujuan", destination, dests, (v) => { destination = v; paint(); })
      ),
      el("div", { class: "mt" }, el("button", { class: "btn", onClick: doSearch }, "Cari jadwal")),
      el("div", { class: "mt" }, loadingSchedules ? spinner() : scheduleList())
    );
  }

  async function doSearch() {
    loadingSchedules = true; paint();
    try { schedules = (await api.schedules(origin, destination)) || []; }
    catch (e) { toast(String(e)); }
    finally { loadingSchedules = false; paint(); }
  }

  api.routes().then((r) => { routes = r || []; loadingRoutes = false; paint(); })
    .catch((e) => { toast(String(e)); loadingRoutes = false; paint(); });

  return { title: "Pesan Tiket", render: () => { paint(); return container; } };
}
