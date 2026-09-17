// Admin home: a left-nav shell with five sections. Each section loads its own
// data. Kept in one file because the sections are small and closely related.
import { api } from "../api.js";
import { el, fmtDate, toast, infoRow, spinner, errBox } from "../ui.js";

// --- Dashboard ---
function stat(label, num) {
  return el("div", { class: "card stat" }, el("div", { class: "num" }, String(num)), el("div", { class: "muted" }, label));
}
function dashboardBody(body) {
  body.innerHTML = ""; body.append(spinner());
  api.adminDashboard().then((s) => {
    body.innerHTML = "";
    body.append(
      el(
        "div",
        { class: "stat-grid" },
        stat("Jadwal hari ini", s.today_schedules),
        stat("Total booking", s.total_bookings),
        stat("Total customer", s.total_customers),
        stat("Tiket terbit", s.tickets_issued)
      )
    );
  }).catch((e) => { body.innerHTML = ""; body.append(errBox(e)); });
}

// --- Generic list (customers / bookings) ---
function listBody(body, loader, renderItem, emptyMsg) {
  body.innerHTML = ""; body.append(spinner());
  loader().then((items) => {
    body.innerHTML = "";
    if (!items || !items.length) { body.append(el("div", { class: "center muted", style: "padding:24px" }, emptyMsg)); return; }
    body.append(...items.map(renderItem));
  }).catch((e) => { body.innerHTML = ""; body.append(errBox(e)); });
}
function renderCustomer(c) {
  return el(
    "div",
    { class: "card" },
    el("strong", {}, c.name),
    el("div", { class: "muted" }, c.phone + (c.email ? " · " + c.email : ""))
  );
}
function renderBooking(b) {
  return el(
    "div",
    { class: "card" },
    el(
      "div",
      { class: "list-row" },
      el(
        "div",
        {},
        el("strong", {}, `${b.booking.booking_code} · ${b.customer.name}`),
        el("div", { class: "muted" }, `${b.origin} → ${b.destination}`),
        el("div", { class: "muted" }, `${fmtDate(b.schedule.departure_at)} · kursi ${b.booking.seat_number}`)
      ),
      el("span", { class: "chip gray" }, b.booking.status)
    )
  );
}

// --- Schedules (list + create form) ---
function fieldWrap(label, input) {
  return el("label", { class: "field grow" }, el("span", {}, label), input);
}
function schedulesBody(body) {
  body.innerHTML = ""; body.append(spinner());
  api.adminSchedules().then((list) => {
    body.innerHTML = "";
    const origin = el("input", { placeholder: "Asal" });
    const dest = el("input", { placeholder: "Tujuan" });
    const cap = el("input", { type: "number", value: "10", min: "1" });
    const dep = el("input", { type: "datetime-local" });

    async function create() {
      if (!origin.value.trim() || !dest.value.trim() || !dep.value || !(+cap.value > 0)) {
        toast("Lengkapi asal, tujuan, kapasitas, dan waktu"); return;
      }
      try {
        await api.adminCreateSchedule({
          origin: origin.value.trim(),
          destination: dest.value.trim(),
          departure_at: new Date(dep.value).toISOString(),
          capacity: +cap.value,
        });
        toast("Jadwal dibuat", true);
        schedulesBody(body);
      } catch (e) { toast(String(e)); }
    }

    body.append(
      el(
        "div",
        { class: "card" },
        el("strong", {}, "Jadwal baru"),
        el("div", { class: "row mt" }, fieldWrap("Asal", origin), fieldWrap("Tujuan", dest)),
        el("div", { class: "row" }, fieldWrap("Kapasitas", cap), fieldWrap("Keberangkatan", dep)),
        el("button", { class: "btn mt", onClick: create }, "Buat jadwal")
      ),
      ...list.map((s) =>
        el(
          "div",
          { class: "card" },
          el(
            "div",
            { class: "list-row" },
            el("div", {}, el("div", {}, `${s.origin} → ${s.destination}`), el("div", { class: "muted" }, fmtDate(s.departure_at))),
            el("div", { class: "end" }, `${s.seats_available}/${s.capacity} kosong`)
          )
        )
      )
    );
  }).catch((e) => { body.innerHTML = ""; body.append(errBox(e)); });
}

// --- Verify a ticket ---
function verifyBody(body) {
  body.innerHTML = "";
  const code = el("input", { placeholder: "mis. TK-DEMO0001" });
  const result = el("div", { class: "mt" });

  async function go() {
    if (!code.value.trim()) return;
    result.innerHTML = ""; result.append(spinner());
    try {
      const r = await api.adminVerify(code.value.trim());
      result.innerHTML = "";
      if (!r) {
        result.append(el("div", { class: "card", style: "background:#ffebee" }, "❌ Tidak valid — tiket tidak ditemukan"));
        return;
      }
      const d = r.detail;
      result.append(
        el("div", { class: "card", style: "background:#e8f5e9" }, "✅ Tiket valid"),
        el(
          "div",
          { class: "card" },
          infoRow("Tiket", d.ticket.ticket_code),
          infoRow("Booking", d.booking.booking_code),
          infoRow("Penumpang", d.customer.name),
          infoRow("Rute", `${d.origin} → ${d.destination}`),
          infoRow("Keberangkatan", fmtDate(d.schedule.departure_at)),
          infoRow("Kursi", String(d.booking.seat_number))
        )
      );
    } catch (e) { result.innerHTML = ""; result.append(errBox(e)); }
  }

  body.append(
    el("label", { class: "field" }, el("span", {}, "Kode tiket atau booking"), code),
    el("button", { class: "btn", onClick: go }, "Verifikasi tiket"),
    result
  );
}

// --- Shell with the left nav ---
export function adminHome() {
  let tab = "dashboard";
  const container = el("div", {});
  const tabs = [
    ["dashboard", "Dashboard"],
    ["customers", "Customer"],
    ["bookings", "Booking"],
    ["schedules", "Jadwal"],
    ["verify", "Verifikasi"],
  ];

  function paint() {
    container.innerHTML = "";
    const nav = el(
      "div",
      { class: "admin-nav" },
      ...tabs.map(([k, label]) =>
        el("button", { class: tab === k ? "active" : "", onClick: () => { tab = k; paint(); } }, label)
      )
    );
    const body = el("div", { class: "admin-body" });
    container.append(el("div", { class: "admin-wrap" }, nav, body));

    if (tab === "dashboard") dashboardBody(body);
    else if (tab === "customers") listBody(body, api.adminCustomers, renderCustomer, "Belum ada customer.");
    else if (tab === "bookings") listBody(body, api.adminBookings, renderBooking, "Belum ada booking.");
    else if (tab === "schedules") schedulesBody(body);
    else if (tab === "verify") verifyBody(body);
  }

  return { title: "Admin", render: () => { paint(); return container; } };
}
