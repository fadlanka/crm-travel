# Travel CRM

A lightweight Travel CRM / customer booking management system built for a technical
assignment. Customers book trips without an account and reschedule their own bookings
using a booking/ticket code; a simple admin panel manages schedules and verifies tickets.

- **Backend:** Go + PostgreSQL, REST API (chi router, pgx driver); QR code digenerate di server
- **Frontend:** HTML + CSS + JavaScript polos (tanpa framework) di `frontend-web/`
- **Frontend lama (arsip):** Flutter Web di `frontend/` — masih ada sebagai referensi, tapi versi yang dipakai adalah vanilla JS
- **Scope:** intentionally small — see `TRAVEL_CRM_CODING_PROMPT.md`

---

## Dokumen penilaian

| Bagian brief | Dokumen |
|--------------|---------|
| Aplikasi + cara menjalankan | README ini |
| Bagian 2 — Arsitektur | [docs/01-arsitektur.md](docs/01-arsitektur.md) |
| Deklarasi AI | [docs/02-deklarasi-ai.md](docs/02-deklarasi-ai.md) |
| Bagian 3 — Analisis Biaya | [docs/03-analisis-biaya.md](docs/03-analisis-biaya.md) |
| Bagian 4 — Riset Penyedia | [docs/04-riset-penyedia.md](docs/04-riset-penyedia.md) |

## Keputusan produk

**Untuk siapa & masalah apa.** CRM untuk **operator travel kecil** dan
**penumpangnya**. Masalah yang dipecahkan: penumpang yang sudah memesan sering harus
menghubungi admin secara manual untuk mengubah jadwal. Solusinya: penumpang
**reschedule sendiri** berbekal kode booking/tiket, tanpa perlu akun.

**Fitur yang dimasukkan (dan kenapa).** Booking tanpa akun, pemilihan rute/jadwal/
kursi, tiket digital + QR, cari & kelola booking, reschedule mandiri (lepas kursi
lama, ambil kursi baru), panel admin (dashboard, customer, booking, jadwal),
verifikasi tiket, data seed + reset. Semua ini adalah inti dari alur "pesan lalu
ubah jadwal" — tak ada yang bisa dibuang tanpa merusak cerita utama.

**Fitur yang sengaja dibuang (dan kenapa).** Pembayaran online, refund, pembatalan,
registrasi akun penumpang, OTP, integrasi WhatsApp/email/SMS, poin loyalti, promo,
manajemen armada/sopir, multi-tenant, analitik lanjutan. Alasan: MVP prototipe
mengasumsikan tiket hanya terbit untuk booking yang sudah lunas, sehingga aturan
reschedule cukup "**tiket ada = boleh reschedule**". Menambah aturan pembayaran/
kebijakan hanya akan mengaburkan keputusan inti yang ingin ditunjukkan.

**Rencana lanjutan jika diteruskan (keterbatasan yang disadari).**
- Manajemen user admin (saat ini satu kredensial bersama, tanpa peran).
- Notifikasi tiket via email (butuh email transaksional — lihat riset penyedia).
- Kebijakan reschedule nyata (batas waktu sebelum keberangkatan, kuota).
- Pagination pada daftar admin (kini memuat semua baris).
- Uji otomatis (unit + integrasi) dan pipeline CI.
- Deploy demo publik + backup terjadwal (lihat analisis biaya untuk targetnya).

---

## Requirements

| Tool | Version used |
|------|--------------|
| Go | 1.27 |
| PostgreSQL | 17 |
| Browser modern | untuk membuka frontend (Chrome/Edge/Firefox) |
| Flutter (opsional) | 3.29 — hanya bila ingin menjalankan frontend lama di `frontend/` |

---

## 1. Database setup

Create a database and user (run once, as the `postgres` superuser):

```sql
CREATE USER travelcrm WITH PASSWORD 'travelcrm';
CREATE DATABASE travel_crm OWNER travelcrm;
```

## 2. Backend

```bash
cd backend
cp .env.example .env        # edit DATABASE_URL if your credentials differ

go run ./cmd/migrate        # create tables (idempotent)
go run ./cmd/seed           # load demo data (resets to a known state)
go run ./cmd/server         # start API on http://localhost:8080
```

Environment variables (`.env`):

```
DATABASE_URL=postgres://travelcrm:travelcrm@localhost:5432/travel_crm?sslmode=disable
PORT=8080
ADMIN_USER=admin
ADMIN_PASSWORD=admin123
```

### Reset to demo state

Re-running the seed truncates every table and reloads the demo data:

```bash
go run ./cmd/seed
```

## 3. Frontend (vanilla HTML/CSS/JS)

Frontend adalah file statis — tidak ada build step. Dua cara menjalankan:

**Cara A — disajikan oleh backend Go (satu origin, paling praktis):**
```bash
cd backend
STATIC_DIR=../frontend-web go run ./cmd/server
# buka http://localhost:8080
```

**Cara B — server statis terpisah:**
```bash
cd frontend-web
python -m http.server 5000   # atau server statis lain
# buka http://localhost:5000  (API tetap di :8080, CORS sudah diaktifkan)
```

Alamat API diatur di `frontend-web/index.html` lewat `window.API_BASE`
(default `http://localhost:8080`). Bila memakai Cara A, boleh dikosongkan
(`""`) agar same-origin.

### Frontend lama (Flutter) — opsional
Versi Flutter masih ada di `frontend/` sebagai arsip:
```bash
cd frontend && flutter run -d chrome
```

---

## Demo data & lookup codes

The seed creates 4 routes, 7 schedules (some today), 4 customers, and 5 bookings.
Handy codes for the demo:

| Booking code | Ticket code | Passenger | Trip |
|--------------|-------------|-----------|------|
| `BK-DEMO01`  | `TK-DEMO0001` | Budi Santoso | Jakarta → Bandung (today) |
| `BK-DEMO04`  | `TK-DEMO0004` | Dewi Lestari | Bandung → Yogyakarta |

Admin login: **admin / admin123**

---

## Trying the full flow

**Customer booking**
1. Open the web app → **Book a trip**
2. Choose origin/destination → **Search schedules** → pick a schedule
3. Choose a seat → enter passenger details → **Confirm booking**
4. The digital ticket (with QR + ticket code) is shown

**Manage & reschedule**
1. **Manage booking** → enter `BK-DEMO01` (or any booking/ticket code) → **Find**
2. **Reschedule this trip** → pick a new schedule → pick a free seat → **Confirm**
3. The old seat is released and the ticket is updated

**Admin**
1. **Admin** → sign in (admin / admin123)
2. Dashboard / Customers / Bookings / Schedules (create new) / Verify a ticket code

---

## API reference

Base URL: `http://localhost:8080`

### Public

| Method | Path | Description |
|--------|------|-------------|
| GET  | `/api/routes` | List routes |
| GET  | `/api/schedules?origin=&destination=` | Upcoming schedules (optional filter) |
| GET  | `/api/schedules/{id}/seats` | Seat map (booked/available) |
| POST | `/api/bookings` | Create a booking (returns ticket) |
| GET  | `/api/bookings/{code}` | Look up by booking **or** ticket code |
| POST | `/api/bookings/{code}/reschedule` | Reschedule a booking |
| GET  | `/api/tickets/{code}` | Ticket lookup |
| GET  | `/api/tickets/{code}/qr` | QR code (PNG) untuk kode tiket/booking |

### Admin (HTTP Basic Auth)

| Method | Path | Description |
|--------|------|-------------|
| GET  | `/api/admin/dashboard` | Summary stats |
| GET  | `/api/admin/customers` | List customers |
| GET  | `/api/admin/bookings` | List bookings |
| GET  | `/api/admin/schedules` | List schedules |
| POST | `/api/admin/schedules` | Create a schedule |
| PUT  | `/api/admin/schedules/{id}` | Update a schedule |
| GET  | `/api/admin/tickets/{code}/verify` | Verify a ticket |

Error responses use standard HTTP codes with a JSON body: `{"error": "..."}`
(`400` invalid input, `404` not found, `409` seat no longer available).

---

## Project structure

```
backend/
  cmd/
    server/    # REST API entrypoint
    migrate/   # applies migrations/*.sql
    seed/      # resets + loads demo data
  internal/
    config/    handler/    middleware/
    db/        model/      repository/    service/
  migrations/  # 001_init.sql
frontend-web/          # frontend aktif — vanilla HTML/CSS/JS
  index.html
  styles.css
  js/
    api.js     # REST client (fetch)
    ui.js      # DOM helper, seat cabin, ticket card, formatting
    app.js     # navigator (screen stack) + semua layar
frontend/              # arsip — versi Flutter lama
  lib/ ...
```

---

## Key design decisions

- **No customer accounts.** Bookings are accessed by booking/ticket code (prompt §3).
- **Reschedule rule is "ticket exists = allowed"** (prompt §2). No payment/OTP/policy checks.
- **Seat exclusivity** is enforced by a partial unique index
  `(schedule_id, seat_number) WHERE status='CONFIRMED'`, so double-booking is impossible
  even under concurrent requests — the DB is the source of truth, not app-level checks.
- **Booking and reschedule run in a single transaction**; on reschedule the old seat is
  freed automatically because the same booking row is moved to the new schedule/seat.
- **Admin auth is HTTP Basic** — the simplest mechanism that satisfies the MVP; no RBAC.
- **Migrations are plain idempotent SQL** applied by a tiny runner (no external CLI).
- **Frontend vanilla JS tanpa framework** — app kecil (~8 layar), setiap baris bisa
  dijelaskan; navigasi memakai screen-stack sederhana (mirip Navigator).
- **QR digenerate di backend Go** — frontend hanya `<img>`; QR hanya berisi kode
  (pointer ke DB), bukan data tiket, jadi tak bisa dipalsukan di sisi klien.

---

## Known limitations (MVP)

- Admin auth is a single shared credential (no user management).
- No pagination on admin lists (fine for demo volumes).
- Ticket is a web page + QR; PDF export is intentionally out of scope (prompt §10).
- `departure_at >= now()` hides past schedules from booking; adjust seed times if needed.
```
