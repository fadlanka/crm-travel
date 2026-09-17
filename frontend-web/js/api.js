// Thin REST client for the Travel CRM API. One method per endpoint.
// No framework — just fetch(). The API base URL comes from window.API_BASE.

const API_BASE = (window.API_BASE || "").replace(/\/$/, "");

// Admin credentials are kept only in memory (cleared on reload / logout).
let adminAuth = null;
export function setAdminAuth(user, pass) {
  adminAuth = "Basic " + btoa(`${user}:${pass}`);
}
export function clearAdminAuth() {
  adminAuth = null;
}

// req performs a fetch and throws an Error (with .status) on non-2xx responses.
async function req(path, { method = "GET", body, admin = false } = {}) {
  const headers = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (admin && adminAuth) headers["Authorization"] = adminAuth;

  const res = await fetch(API_BASE + path, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  const text = await res.text();
  const data = text ? JSON.parse(text) : null;

  if (!res.ok) {
    const err = new Error((data && data.error) || `Request gagal (${res.status})`);
    err.status = res.status;
    throw err;
  }
  return data;
}

export const api = {
  // URL for the server-generated QR image of a ticket/booking code.
  qrUrl: (code) => `${API_BASE}/api/tickets/${encodeURIComponent(code)}/qr`,

  // --- Customer / booking ---
  routes: () => req("/api/routes"),

  schedules: (origin, destination) => {
    const q = new URLSearchParams();
    if (origin) q.set("origin", origin);
    if (destination) q.set("destination", destination);
    const s = q.toString();
    return req("/api/schedules" + (s ? `?${s}` : ""));
  },

  seats: (scheduleId) => req(`/api/schedules/${scheduleId}/seats`),

  createBooking: (payload) => req("/api/bookings", { method: "POST", body: payload }),

  getBooking: (code) => req(`/api/bookings/${encodeURIComponent(code)}`),

  reschedule: (code, payload) =>
    req(`/api/bookings/${encodeURIComponent(code)}/reschedule`, { method: "POST", body: payload }),

  // --- Admin (Basic Auth) ---
  adminDashboard: () => req("/api/admin/dashboard", { admin: true }),
  adminCustomers: () => req("/api/admin/customers", { admin: true }),
  adminBookings: () => req("/api/admin/bookings", { admin: true }),
  adminSchedules: () => req("/api/admin/schedules", { admin: true }),
  adminCreateSchedule: (payload) =>
    req("/api/admin/schedules", { method: "POST", body: payload, admin: true }),

  // Returns { valid, detail } when found, or null when the ticket does not exist.
  adminVerify: async (code) => {
    try {
      return await req(`/api/admin/tickets/${encodeURIComponent(code)}/verify`, { admin: true });
    } catch (e) {
      if (e.status === 404) return null;
      throw e;
    }
  },
};
