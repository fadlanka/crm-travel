# Riset & Pemilihan Penyedia — Travel CRM

> **Wajib diverifikasi:** harga dikumpulkan sekitar **September 2026** dari halaman
> resmi & agregator (tautan di tiap bagian). Konfirmasi ulang di halaman resmi
> sebelum submit. Kurs USD 1 = Rp 16.500 (asumsi). Nama penyedia & alasan di sini
> adalah **rekomendasi** untuk kamu setujui/ubah — kamu yang mempertahankannya nanti.

---

## Tabel ringkas pilihan (bisa dibaca satu menit)

| Komponen | Penyedia dipilih | Biaya/bln | Alasan memilih (1 kalimat) | Juga dipertimbangkan | Kenapa tidak dipilih |
|----------|------------------|----------:|----------------------------|----------------------|----------------------|
| **Compute** | Biznet Gio NEO Lite (Jakarta) | Rp 50–95k | Latensi Jakarta terendah + bandwidth gratis tanpa kuota + bayar Rupiah | IDCloudHost; Render | IDCloudHost setara (cadangan); Render tak punya region Indonesia & egress metered |
| **Database** | PostgreSQL self-host di VPS | Rp 0 | Gratis marginal, data tetap di Indonesia, lock-in terendah (pg_dump) | Neon; Supabase | Neon/Supabase = data lintas negara + free tier Supabase pause saat idle |
| **Object storage** | *Tidak dipakai* (R2 bila perlu) | Rp 0 | Aplikasi tak menyimpan berkas; QR di-generate di klien | Backblaze B2; AWS S3 | Bila perlu, R2 dipilih karena **egress $0**; S3 egress mahal |
| **Email transaksional** | *Tidak dipakai* (Resend bila perlu) | Rp 0 | Kirim email di luar ruang lingkup MVP | Amazon SES; Brevo | Bila perlu, Resend paling simpel; SES untuk volume besar |
| **Autentikasi** | Bawaan (Basic Auth + kode tiket) | Rp 0 | Cukup untuk MVP, tanpa dependensi | Auth0; Supabase Auth | Berlebihan untuk satu kredensial admin |
| **Monitoring** | UptimeRobot (free) | Rp 0 | Free 50 monitor + status page | Better Stack; Pingdom | Free tier lebih kecil |
| **Error tracking** | Sentry (free Developer) | Rp 0 | SDK Go + JS, free 5.000 error/bln | GlitchTip; self-host | Butuh operasional tambahan |
| **Domain** | Cloudflare Registrar `.com` | Rp 14k | Harga at-cost tanpa markup | Niagahoster/Hostinger `.id` | `.id` perpanjangan lebih mahal |
| **SSL** | Let's Encrypt (Caddy) | Rp 0 | Otomatis, gratis, standar industri | Cloudflare Universal SSL | Setara; keduanya gratis |
| **CDN** | Cloudflare (free) | Rp 0 | PoP Jakarta, bandwidth cacheable gratis | bunny.net; Fastly | Berbayar sejak awal |
| **Backup** | pg_dump → R2/Backblaze B2 | Rp 0–50k | Murah, off-site, portabel | Snapshot VPS Biznet | Snapshot terkunci di penyedia |

---

## Perbandingan mendalam (3 penyedia) untuk 4 komponen penentu

### 1. Compute

| | Biznet Gio NEO Lite **(dipilih)** | IDCloudHost Cloud VPS | Render |
|---|---|---|---|
| Region ke Indonesia | **Jakarta** | Jakarta + Singapura | US / EU / Singapura |
| Harga awal | Rp 50.000 (1 vCPU/1 GB/60 GB) | Rp 50.000 (1 CPU/1 GB/20 GB); eXtreme dari Rp 87.000 | Free (spin-down 15 mnt); Starter $7 |
| Bandwidth/egress | **Gratis tanpa kuota (10 Gbps)** | Gratis / free bandwidth | Metered di luar kuota plan |
| Pembayaran | Rupiah (transfer/VA) | Rupiah (per jam–bulan) | Kartu kredit USD |
| Operasional | Manual (OS, SSL, update) | Manual | **Zero-ops (PaaS)** |

**Kenapa Biznet layak dipakai:** untuk pengguna Indonesia, latensi Jakarta ~5–15 ms
mengalahkan Singapura (~20–40 ms) dan US (150 ms+). Bandwidth gratis tanpa kuota
menghapus pos egress yang di skala B jadi termahal. Pembayaran Rupiah menghilangkan
kebutuhan kartu kredit internasional. vCPU/RAM dedicated (bukan shared burst).

**Kenapa dua lainnya tidak dipilih:**
- **IDCloudHost** sangat mirip dan jadi **cadangan utama** — juga Jakarta, juga
  bandwidth gratis, model bayar per jam yang fleksibel. Biznet dipilih karena jaminan
  bandwidth 10 Gbps tanpa kuota dinyatakan eksplisit; keduanya bisa ditukar tanpa
  ubah kode (Docker Compose).
- **Render** menang di kenyamanan (deploy dari Git, zero-ops) tetapi kalah di konteks
  Indonesia: tak ada region Indonesia (data keluar negeri → gesekan UU PDP), egress
  metered, bayar kartu USD, dan free tier **spin-down** setelah 15 menit idle =
  cold start yang buruk untuk demo yang ditinggal reviewer. Render juga
  [dilaporkan makin mahal](https://servercompass.app/blog/render-pricing-is-it-worth-it)
  setelah perubahan paket April 2026.

Sumber: [Biznet NEO Lite](https://www.biznetgio.com/product/neo-lite) ·
[IDCloudHost pricing](https://idcloudhost.com/pricing/) ·
[Render pricing](https://render.com/pricing)

### 2. Database

| | PostgreSQL self-host **(dipilih)** | Neon | Supabase |
|---|---|---|---|
| Model | Di VPS sendiri | Serverless managed PG | Managed PG + Auth/Realtime |
| Free tier | — (termasuk compute) | 0,5 GB, 100 CU-jam/bln, scale-to-zero | 500 MB, 1 GB file; **pause 1 mgg idle** |
| Harga naik | Rp 0 marginal | Launch $0,106/CU-jam + $0,35/GB | Pro $25/bln/proyek |
| Region | **Jakarta (di VPS)** | US/EU (managed) | termasuk Singapura |
| Lock-in | **Rendah (pg_dump)** | Sedang (fitur branching) | Sedang–tinggi (bundel fitur) |

**Kenapa self-host layak:** PostgreSQL adalah PostgreSQL — pindah gampang lewat
`pg_dump`. Data tetap di Indonesia (UU PDP paling mudah dipatuhi). Biaya marginal
Rp0 karena berjalan di VPS yang sudah dibayar. Kontrol penuh atas constraint &
transaksi yang jadi inti sistem (lihat partial unique index di arsitektur).

**Kenapa dua lainnya tidak dipilih (untuk skala A):**
- **Neon** unggul untuk zero-ops & scale-to-zero, dan free tier-nya cukup untuk A.
  Tak dipilih karena menambah dependensi managed di US (data lintas negara) padahal
  DB kita kecil dan sudah punya VPS. **Dipertimbangkan ulang di skala B** bila butuh
  HA/PITR tanpa mengelola sendiri. (Storage turun ~80% ke $0,35/GB setelah Databricks
  mengakuisisi Neon, 2025.)
- **Supabase** membundel Auth/Realtime yang tak kita butuh, dan **proyek free-nya
  pause setelah 1 minggu idle** — buruk untuk demo penilaian yang bisa menganggur.
  Pro melompat ke $25/bln.

Sumber: [Neon pricing](https://neon.com/pricing) ·
[Supabase pricing](https://supabase.com/pricing)

### 3. Object storage — *saat ini tidak dipakai*

Aplikasi tidak menyimpan berkas (QR & tiket di-render di klien). Tetap dibandingkan
seandainya fitur unggah/PDF ditambah:

| | Cloudflare R2 **(pilihan bila perlu)** | Backblaze B2 | AWS S3 Standard |
|---|---|---|---|
| Storage | $0,015/GB | **$0,00695/GB** (termurah) | $0,023/GB |
| **Egress** | **$0 (gratis)** | Gratis 3× rata2 storage, lalu $0,01/GB; gratis via Cloudflare | **$0,09/GB** |
| Free tier | 10 GB + 1 jt Class A ops | 10 GB | 100 GB egress/bln (12 bln) |
| Operasi | Class A $4,50/jt, Class B $0,36/jt | murah | berlapis |

**Kenapa R2 (bila perlu):** egress $0 berpasangan sempurna dengan CDN Cloudflare —
tak ada biaya kejut saat file sering diunduh. S3-compatible. **Backblaze B2** lebih
murah untuk storage murni (arsip backup) dan egress-nya gratis ke Cloudflare — itu
sebabnya B2 kami pakai untuk **backup** di skala B. **AWS S3** paling mahal karena
egress $0,09/GB dan paling kompleks; keunggulannya hanya region Singapura + ekosistem.

Sumber: [R2 pricing](https://developers.cloudflare.com/r2/pricing) ·
[Backblaze B2](https://www.backblaze.com/cloud-storage/pricing) ·
[S3 pricing](https://aws.amazon.com/s3/pricing/)

### 4. Email transaksional — *saat ini tidak dipakai*

Pengiriman email di luar ruang lingkup MVP. Dibandingkan seandainya reminder/tiket
via email ditambah:

| | Resend **(pilihan bila perlu)** | Amazon SES | Brevo |
|---|---|---|---|
| Free tier | 3.000/bln (**100/hari**) | 3.000/bln (12 bln pertama) | 300/hari |
| Harga | Pro $20/50k, Scale $90/100k | **$0,10/1.000** (termurah skala) | $9/5k |
| Setup | **Paling simpel (API)** | Verifikasi domain + keluar sandbox | Menengah |
| Region | global | termasuk Singapura | global |

**Kenapa Resend (bila perlu):** API paling ramah developer, free tier paling murah
hati untuk volume kecil. Kekurangan: cap 100/hari di free, perusahaan relatif baru.
**Amazon SES** termurah untuk volume besar ($5 untuk 50.000 email di skala B) — jadi
pilihan bila email benar-benar dipakai di B — tetapi setup lebih repot. **Brevo**
berorientasi marketing dengan plafon free lebih rendah untuk kebutuhan transaksional.

Sumber: [Resend pricing](https://resend.com/pricing) ·
[Amazon SES pricing](https://aws.amazon.com/ses/pricing/) ·
[Brevo pricing](https://www.brevo.com/pricing/)

---

## Komponen sisanya (1–2 kalimat)

- **Autentikasi.** Memakai mekanisme bawaan: HTTP Basic Auth untuk admin + akses via
  kode booking/tiket untuk penumpang; tidak memakai layanan auth pihak ketiga.
  Konsekuensi: tak ada SSO/MFA dan satu kredensial admin bersama — diterima untuk MVP.
- **Monitoring/uptime.** UptimeRobot free (50 monitor, cek 5 menit) untuk memantau
  endpoint `/health` + halaman status; Rp0 dan tak akan menyentuh batas free.
- **Error tracking.** Sentry free Developer (5.000 error/bln, retensi 30 hari), SDK di
  Go dan JavaScript; naik ke Team $26/bln bila terlampaui (lihat titik pindah).
- **Domain.** Cloudflare Registrar `.com` at-cost ($10,44/th, naik $11,15 per 1 Nov
  2026); alternatif `.id` lewat registrar lokal (~Rp230k/th) bila ingin identitas
  Indonesia — konsekuensinya perpanjangan `.id` lebih mahal daripada promo awal.
- **SSL.** Let's Encrypt via Caddy (perpanjang otomatis) atau Cloudflare Universal
  SSL — keduanya gratis; memakai fasilitas bawaan penyedia yang sudah dipilih.
- **CDN.** Cloudflare free, memanfaatkan PoP Jakarta — bandwidth cacheable gratis,
  kemenangan latensi besar untuk pengguna Indonesia; konsekuensi ToS: tak untuk media
  besar non-HTML (tak relevan, aset kami web app).
- **Backup.** `pg_dump` terjadwal setiap malam → Cloudflare R2 (free 10 GB) di skala
  A, Backblaze B2 ($6,95/TB) di skala B; disimpan off-VPS. Konsekuensi: RPO 24 jam,
  diterima untuk MVP.

---

## Konteks Indonesia (wajib dibahas)

- **Lokasi server & latensi.** Biznet Gio & IDCloudHost punya **region Jakarta**
  (~5–15 ms ke pengguna Indonesia). Cloudflare punya **PoP Jakarta**. Penyedia global
  (Render/Neon/Supabase/AWS) terdekat = **Singapura** (~20–40 ms) atau US (150 ms+).
  Stack pilihan menahan compute + data + CDN di dalam/dekat Jakarta.
- **Cara pembayaran.** Biznet/IDCloudHost menerima **Rupiah** (transfer bank, Virtual
  Account, e-wallet) — tanpa kartu kredit internasional. Render/Neon/AWS/Cloudflare
  butuh kartu kredit/debit USD. **Bila pembayaran gagal di jatuh tempo:** penyedia
  lokal umumnya memberi masa tenggang lalu suspend (data ditahan sementara); domain
  Cloudflare perlu kartu valid untuk auto-renew — kegagalan bisa berujung domain
  lapse, risiko yang harus dipantau.
- **UU Perlindungan Data Pribadi (UU 27/2022).** CRM menyimpan data pribadi (nama,
  telepon, email). Menyimpan data di **VPS Jakarta (self-host)** menjaga data
  di dalam negeri → kepatuhan paling sederhana. Memakai DB managed di US (Neon)
  berarti **transfer data lintas negara**, yang menuntut dasar hukum/persetujuan &
  pengamanan tambahan. Ini salah satu alasan kuat memilih VPS lokal + self-host.
- **Penyedia lokal.** Dibandingkan **Biznet Gio** dan **IDCloudHost** (keduanya cocok
  — Biznet dipilih karena jaminan bandwidth gratis). **DomaiNesia** kuat untuk
  shared/managed hosting tetapi VPS-nya kurang fleksibel/lebih mahal untuk kebutuhan
  ini. **Kesimpulan: untuk kasus ini penyedia lokal justru paling cocok** — berbeda
  dari banyak web app yang defaultnya langsung ke PaaS global.

---

## Yang membuat perbandingan ini serius (checklist brief)

- **Harga dari halaman resmi** — tautan resmi disertakan tiap bagian; tandai
  September 2026, wajib diverifikasi ulang saat submit.
- **Batas free tier + yang terjadi saat terlampaui** — disebut per komponen
  (mis. Sentry 5.000 error → $26/bln; Resend 100/hari; Neon 100 CU-jam).
- **Hal yang tidak enak juga disebut** — Render makin mahal + spin-down; Fly.io
  [mematikan free tier permanen](https://expresstech.io/7-fly-io-alternatives-in-2026-real-pricing-after-the-free-tier-died/);
  Supabase free pause saat idle; overage Sentry melompat; harga `.com` naik 1 Nov
  2026 (Verisign); ToS media Cloudflare free.
- **Biaya keluar (exit cost)** — dibahas di [`03-analisis-biaya.md`](03-analisis-biaya.md):
  self-host + Docker Compose = lock-in terendah; managed = ekspor lebih repot.
- **Egress** — dihitung eksplisit; jadi pos terbesar di skala B pada penyedia global,
  Rp0 pada VPS lokal berbandwidth gratis.
