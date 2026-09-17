// Small UI helpers: a tiny DOM builder, formatting, toast, and the two
// reusable widgets (seat cabin + ticket card). No framework.
import { api } from "./api.js";

// el("div", {class:"card", onClick: fn}, child1, child2, ...) -> HTMLElement
export function el(tag, props, ...children) {
  const node = document.createElement(tag);
  if (props) {
    for (const [k, v] of Object.entries(props)) {
      if (v == null || v === false) continue;
      if (k === "class") node.className = v;
      else if (k.startsWith("on") && typeof v === "function") {
        node.addEventListener(k.slice(2).toLowerCase(), v);
      } else if (v === true) node.setAttribute(k, "");
      else node.setAttribute(k, v);
    }
  }
  for (const child of children.flat()) {
    if (child == null || child === false) continue;
    node.append(child.nodeType ? child : document.createTextNode(String(child)));
  }
  return node;
}

// Format an ISO departure timestamp for display, e.g. "Kam, 17 Sep 2026, 15.00".
const dateFmt = new Intl.DateTimeFormat("id-ID", {
  weekday: "short",
  day: "numeric",
  month: "short",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});
export function fmtDate(iso) {
  return dateFmt.format(new Date(iso));
}

// Toast notification. ok=true renders green, otherwise red (error).
let toastTimer = null;
export function toast(message, ok = false) {
  const t = document.getElementById("toast");
  t.textContent = message;
  t.className = "toast" + (ok ? " ok" : "");
  t.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => (t.hidden = true), 4000);
}

export function infoRow(label, value) {
  return el("div", { class: "info" }, el("div", { class: "k" }, label), el("div", { class: "v" }, value));
}

export function spinner() {
  return el("div", { class: "spinner" }, "Memuat…");
}

// Small red error box, shown when a screen fails to load data.
export function errBox(e) {
  return el("div", { class: "center", style: "color:var(--danger);padding:16px" }, "Error: " + e);
}

// --- Seat cabin ---
// Draws the seat map as a real vehicle cabin (right-hand drive):
// seat 1 up front next to the driver, then rows of three behind
// (one on the left, two on the right, split by an aisle).
export function seatCabin(seatMap, selectedSeat, onSelect) {
  const booked = new Set(seatMap.booked);
  const capacity = seatMap.capacity;

  function seatCell(seat) {
    if (seat == null) return el("div", { class: "seat empty" });
    const isBooked = booked.has(seat);
    let cls = "seat";
    if (isBooked) cls += " booked";
    else if (seat === selectedSeat) cls += " selected";
    return el(
      "button",
      {
        class: cls,
        onClick: isBooked ? null : () => onSelect(seat),
        disabled: isBooked,
      },
      el("span", { class: "ico" }, "\u{1F4BA}"), // seat emoji
      el("span", {}, String(seat))
    );
  }

  const driverCell = el("div", { class: "seat driver", title: "Supir" }, el("span", { class: "ico" }, "\u{1F535}"));

  const cabin = el("div", { class: "cabin" }, el("div", { class: "cabin-label" }, "DEPAN"));

  // Front row: seat 1 (left) | aisle | empty | driver (right = right-hand drive).
  cabin.append(
    el(
      "div",
      { class: "seat-row" },
      seatCell(1), // kapasitas dijamin > 0 oleh CHECK di database
      el("div", { class: "aisle" }),
      el("div", { class: "seat empty" }),
      driverCell
    ),
    el("hr")
  );

  // Rear rows: seats 2..capacity, three per row [left | aisle | inner | window].
  const rear = [];
  for (let s = 2; s <= capacity; s++) rear.push(s);
  for (let i = 0; i < rear.length; i += 3) {
    const row = rear.slice(i, i + 3);
    cabin.append(
      el(
        "div",
        { class: "seat-row" },
        seatCell(row[0]),
        el("div", { class: "aisle" }),
        seatCell(row[1]),
        seatCell(row[2])
      )
    );
  }

  return el(
    "div",
    {},
    el(
      "div",
      { class: "legend" },
      legendItem("var(--primary)", "Dipilih"),
      legendItem("#fff", "Tersedia"),
      legendItem("#b7bec7", "Terisi"),
      legendItem("#37474f", "Supir")
    ),
    cabin
  );
}

function legendItem(color, label) {
  const sw = el("span", { class: "swatch" });
  sw.style.background = color;
  return el("div", { class: "item" }, sw, label);
}

// --- Ticket card (with server-generated QR) ---
export function ticketCard(d) {
  const qr = el("img", {
    src: api.qrUrl(d.ticket.ticket_code),
    alt: "QR " + d.ticket.ticket_code,
    width: 160,
    height: 160,
  });
  return el(
    "div",
    { class: "card" },
    el(
      "div",
      { class: "row between" },
      el("strong", {}, "Tiket Digital"),
      el("span", { class: "chip ok" }, d.ticket.status)
    ),
    el("div", { class: "center mt" }, qr),
    el("div", { class: "center", style: "font-size:18px;font-weight:700;letter-spacing:2px" }, d.ticket.ticket_code),
    el("hr"),
    infoRow("Kode booking", d.booking.booking_code),
    infoRow("Penumpang", d.customer.name),
    infoRow("Telepon", d.customer.phone),
    infoRow("Rute", `${d.origin} → ${d.destination}`),
    infoRow("Keberangkatan", fmtDate(d.schedule.departure_at)),
    infoRow("Kursi", String(d.booking.seat_number)),
    infoRow("Status", d.booking.status)
  );
}
