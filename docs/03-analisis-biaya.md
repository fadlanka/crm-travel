# Analisis Biaya — Travel CRM

Biaya menjalankan sistem untuk **satu bulan** dan **satu tahun**, pada dua skenario.

> **Yang perlu kamu verifikasi sebelum submit (brief mewajibkan harga dari halaman
> resmi):** angka di bawah dikumpulkan sekitar **September 2026** dari halaman harga
> penyedia dan agregator. Tautan resmi ada di [`04-riset-penyedia.md`](04-riset-penyedia.md).
> Konfirmasi ulang tiap angka di halaman resmi, dan sesuaikan bila berbeda.
>
> **Kurs:** USD 1 = **Rp 16.500** (asumsi — cek kurs saat submit).

---

## Asumsi beban (ditulis eksplisit — angka tanpa asumsi dianggap tidak ada)

| Parameter | Skenario A (skala awal) | Skenario B (×100) |
|-----------|-------------------------|-------------------|
| Operator admin | 1 | ~5 |
| Booking / bulan | 500 | 50.000 |
| Total customer (record) | 3.000 | 300.000 |
| Total booking / tahun | 6.000 | 600.000 |
| Jadwal aktif | ~300 | ~30.000 |
| Request API / bulan | ~50.000 | ~5.000.000 |
| Ukuran database | ~1 GB | ~75–100 GB |
| Berkas tersimpan (media) | 0 | 0 |
| Egress (data keluar) / bulan | ~30 GB | ~3.000 GB (3 TB) |
| Email transaksional / bulan | 0 | 0 |

**Catatan penting dua asumsi bernilai nol:**
- **Berkas = 0.** Aplikasi tidak menyimpan file. Tiket digital + QR di-*generate* di sisi
  klien (browser), tidak ada upload/PDF. Maka **object storage tidak dipakai**.
- **Email = 0.** Pengiriman email/SMS eksplisit **di luar ruang lingkup MVP**
  (lihat spec §6). Maka **email transaksional tidak dipakai**.

Keduanya tetap muncul di tabel biaya (diisi "tidak dipakai" + alasan + berapa
biayanya seandainya ditambahkan), sesuai aturan brief.

---

## Stack yang dihitung (jalur yang direkomendasikan)

Satu **VPS di Jakarta** menjalankan biner Go + PostgreSQL + Caddy (auto-SSL) via
Docker Compose, dengan **Cloudflare (free)** sebagai DNS/CDN/SSL di depan. Alasan
lengkap ada di [`04-riset-penyedia.md`](04-riset-penyedia.md). Keunggulan yang
langsung berdampak ke biaya: **VPS lokal ini punya bandwidth gratis tanpa kuota**,
sehingga **egress praktis Rp0** — pos yang justru paling mahal di penyedia global.

---

## Skenario A — rincian per komponen

| Komponen | Penyedia | Biaya/bulan | Biaya/tahun | Catatan free-tier / batas |
|----------|----------|------------:|------------:|---------------------------|
| **Compute** | Biznet Gio NEO Lite (2 vCPU/2 GB, verif.) | Rp 95.000 | Rp 1.140.000 | Tier 1 vCPU/1 GB = Rp 50.000 bila RAM cukup |
| **Database** | PostgreSQL self-host (di VPS) | Rp 0 | Rp 0 | Termasuk dalam compute; alternatif Neon free 0,5 GB/100 CU-jam |
| **Object storage** | — *tidak dipakai* | Rp 0 | Rp 0 | Bila ditambah: Cloudflare R2 free 10 GB, lalu $0,015/GB |
| **Email transaksional** | — *tidak dipakai* | Rp 0 | Rp 0 | Bila ditambah: Resend free 3.000/bln (100/hari) |
| **Autentikasi** | Bawaan (Basic Auth admin) | Rp 0 | Rp 0 | Bukan layanan pihak ketiga |
| **Monitoring/uptime** | UptimeRobot (free) | Rp 0 | Rp 0 | Free 50 monitor, cek 5 menit |
| **Error tracking** | Sentry (free Developer) | Rp 0 | Rp 0 | Free 5.000 error/bln, retensi 30 hari |
| **Domain** | Cloudflare Registrar `.com` (at-cost) | Rp 14.350 | Rp 172.000 | $10,44/th; naik ke $11,15 per 1 Nov 2026 (Verisign) |
| **SSL** | Let's Encrypt (Caddy) / Cloudflare | Rp 0 | Rp 0 | Perpanjang otomatis |
| **CDN** | Cloudflare (free) | Rp 0 | Rp 0 | Bandwidth cacheable gratis, PoP Jakarta |
| **Backup** | pg_dump → Cloudflare R2 (free 10 GB) | Rp 0 | Rp 0 | RPO 24 jam; dump < 10 GB |
| **Egress** | Bandwidth VPS (gratis tanpa kuota) | Rp 0 | Rp 0 | Di penyedia global ~30 GB = ~$2,7 |
| **TOTAL** | | **± Rp 109.350** | **± Rp 1.312.000** | (~$80/tahun) |

---

## Skenario B — rincian per komponen (beban ×100)

Pada ×100, **arsitektur diubah**: PostgreSQL dipindah ke VPS-nya sendiri, terpisah
dari aplikasi (titik pindah dijelaskan di bawah).

| Komponen | Penyedia | Biaya/bulan | Biaya/tahun | Catatan |
|----------|----------|------------:|------------:|---------|
| **Compute (app)** | Biznet NEO Lite 4 vCPU/8 GB | Rp 245.000 | Rp 2.940.000 | |
| **Compute (DB)** | Biznet NEO Lite 4 vCPU/8 GB | Rp 245.000 | Rp 2.940.000 | DB dipisah dari app |
| **Database** | PostgreSQL self-host (di VPS DB) | Rp 0 | Rp 0 | Termasuk compute DB; alt. Neon Scale ~$26+/bln |
| **Object storage** | — *tidak dipakai* | Rp 0 | Rp 0 | Bila ditambah: R2 $0,015/GB, egress $0 |
| **Email transaksional** | — *tidak dipakai* | Rp 0 | Rp 0 | Bila reminder ditambah: SES 50.000 email = ~$5 = Rp 82.500 |
| **Autentikasi** | Bawaan | Rp 0 | Rp 0 | |
| **Monitoring/uptime** | UptimeRobot (free) | Rp 0 | Rp 0 | 50 monitor masih cukup |
| **Error tracking** | Sentry (free → Team bila >5.000 error) | Rp 0–429.000 | Rp 0–5.148.000 | Team $26/bln bila error > 5.000/bln |
| **Domain** | Cloudflare `.com` | Rp 14.350 | Rp 172.000 | **Tidak naik** dengan skala |
| **SSL** | Let's Encrypt / Cloudflare | Rp 0 | Rp 0 | |
| **CDN** | Cloudflare (free) | Rp 0 | Rp 0 | Pro $25/bln hanya bila butuh WAF/aturan lanjutan |
| **Backup** | pg_dump → Backblaze B2 (~200–300 GB) | Rp 49.500 | Rp 594.000 | ~$3/bln ($6,95/TB) |
| **Egress** | Bandwidth VPS (gratis) | Rp 0 | Rp 0 | **Di AWS: 3 TB × $0,09 = ~Rp 4.455.000/bln** |
| **TOTAL (dasar)** | | **± Rp 553.850** | **± Rp 6.646.000** | Sentry masih free |
| **TOTAL (bila Sentry Team)** | | **± Rp 982.850** | **± Rp 11.794.000** | |

---

## Bacaan yang diminta brief: biaya tidak naik lurus

Dari A ke B beban naik **100×**, tetapi biaya dasar hanya naik dari ~Rp 109k ke
~Rp 554k — sekitar **5×**. Kenapa tiap komponen berperilaku beda:

| Perilaku | Komponen | Penjelasan |
|----------|----------|------------|
| **Tetap (0× naik)** | Domain, SSL | Harga domain tak peduli jumlah pengguna |
| **Tetap dalam batas gratis** | CDN, Monitoring, (Error tracking di A) | Free tier menutup A **dan** B |
| **Melompat saat batas terlampaui** | Error tracking | Gratis sampai 5.000 error/bln, lalu lompat ke $26/bln |
| **Naik bertahap (bukan 100×)** | Compute | Dari 1 VPS kecil → 2 VPS lebih besar, bukan 100 mesin |
| **Nyaris nol karena pilihan arsitektur** | Egress | VPS lokal = bandwidth gratis; di cloud global egress justru jadi pos **terbesar** di B |
| **Arsitektur harus diubah sebelum skala** | Database | Di B, PostgreSQL dipisah ke VPS sendiri (bukan sekadar naik paket) |

**Pelajaran utama (yang paling penting untuk dipertahankan):** komponen yang di
penyedia global akan meledak di skala B adalah **egress** — 3 TB/bulan ≈
Rp 4,45 juta/bulan di AWS, lebih besar dari seluruh tagihan kami digabung. Pilihan
VPS lokal berbandwidth gratis membuat pos itu **Rp0**. Inilah alasan biaya B tetap
masuk akal.

---

## Titik pindah (kapan harus naik kelas / ganti paket / ubah arsitektur)

| Komponen | Titik pindah | Aksi |
|----------|--------------|------|
| **Compute** | DB working-set > RAM, atau CPU app jenuh (~10–20k booking/bln) | Pisahkan DB ke VPS sendiri, naikkan RAM |
| **Database** | Butuh HA, PITR, atau replika baca; DB > RAM | Pindah ke managed (Neon/RDS) atau tambah replika |
| **Error tracking** | > 5.000 error/bln | Sentry Team $26/bln, atau turunkan noise |
| **CDN (Cloudflare free)** | Butuh WAF, rate-limit, aturan lanjutan | Cloudflare Pro $25/bln |
| **Backup** | Dump > 10 GB (lewat R2 free) | Backblaze B2 $6,95/TB |
| **Egress** | **Hanya** jika pindah dari VPS free-bandwidth ke cloud metered | Hitung ulang — bisa jadi pos terbesar |
| **Domain `.com`** | 1 Nov 2026 | Harga wholesale naik $10,26 → $10,97 (Verisign) |

---

## Ongkos kalau ternyata salah pilih (exit cost / lock-in)

Stack ini sengaja dipilih untuk **lock-in serendah mungkin**:

- **Database:** PostgreSQL self-host → `pg_dump` menghasilkan berkas portabel; pindah
  ke VPS/penyedia lain hitungan jam. Bandingkan managed proprietary (fitur khusus,
  ekspor lebih repot).
- **Compute:** Docker Compose portabel — jalan di VPS mana pun. Pindah dari Biznet ke
  IDCloudHost/DomaiNesia tanpa ubah kode.
- **CDN/DNS:** Cloudflare free — keluar = ganti nameserver, tak ada data terkunci.
- **Domain:** registrar at-cost, transfer keluar standar (biaya ~1 tahun perpanjangan).

Konsekuensinya (harga dari kemudahan pindah): kami menanggung sendiri operasional —
backup, patch OS, upgrade PostgreSQL, dan tidak ada HA otomatis. Untuk skala A–B ini
pertukaran yang disengaja; disebut juga sebagai keterbatasan di
[`01-arsitektur.md`](01-arsitektur.md) dan README.
