# Deklarasi Penggunaan AI

Sesuai syarat penugasan (Bagian "Penggunaan AI"), berkas ini menerangkan di
bagian mana AI dipakai dan sejauh apa. Tidak ada nilai yang dikurangi karena
memakai AI; yang dinilai adalah kejujuran deklarasi dan kemampuan menjelaskan
kode sendiri.

> **Catatan untuk dilengkapi pemilik proyek:** sesuaikan bagian yang ditandai
> _[isi sendiri]_ agar mencerminkan porsi keterlibatanmu yang sebenarnya.
> Deklarasi ini harus jujur karena kamu akan diminta menjelaskan kode ini,
> termasuk bagian yang dihasilkan AI.

## Alat yang dipakai

- **Claude Code** (AI coding assistant) — model Claude Opus.

## Sejauh apa AI dipakai

| Bagian | Peran AI | Peran manusia |
|--------|----------|---------------|
| Penentuan ruang lingkup & keputusan produk | — | **Manusia.** Ditulis lebih dulu di `TRAVEL_CRM_CODING_PROMPT.md` (untuk siapa, fitur masuk/dibuang, aturan reschedule) |
| Skema database & migrasi | Draf SQL dari model yang diminta | Review & keputusan constraint |
| Backend Go (handler/service/repository) | Generate sebagian besar kode implementasi | Review, uji, penjelasan |
| Frontend web (vanilla HTML/CSS/JS) | Generate layar & alur dari daftar UI di spec | Review, uji, penjelasan |
| Frontend Flutter (arsip di `frontend/`) | Versi awal, kemudian diganti vanilla JS | Keputusan mengganti stack |
| Data seed & runner migrasi | Generate | Review |
| Dokumen arsitektur & README | Draf dari kode yang sudah ada | Verifikasi kebenaran |
| Analisis biaya & riset penyedia | _[isi sendiri: riset harga & pemilihan penyedia sebaiknya kamu lakukan/verifikasi sendiri — angka harus dari halaman resmi]_ | **Manusia** |

## Cara kerja

Spesifikasi (`TRAVEL_CRM_CODING_PROMPT.md`) ditulis manusia lebih dulu sebagai
kontrak produk. AI diminta mengimplementasikan sesuai spec itu, dengan aturan
eksplisit untuk tidak menambah fitur atau mengubah aturan bisnis secara diam-diam.
Setiap milestone menghasilkan kode yang bisa dijalankan, lalu ditinjau.

## Yang dipahami penulis

Titik-titik yang paling mungkin ditanya di sesi lanjutan dan perlu dikuasai:

- **Kenapa keunikan kursi memakai partial unique index**, bukan cek di aplikasi —
  dan apa yang terjadi kalau dua orang memesan kursi sama bersamaan.
- **Kenapa booking & reschedule dibungkus transaksi**, dan apa yang terjadi bila
  gagal di tengah (rollback).
- **Kenapa reschedule memindahkan baris booking** alih-alih membuat booking baru,
  dan bagaimana kursi lama menjadi bebas.
- **Kenapa penumpang tanpa login** dan apa konsekuensi keamanannya.
- **Kenapa Go + vanilla JS + PostgreSQL** (dan kenapa Flutter diganti), dan dalam kondisi apa akan memilih lain
  (lihat tabel keputusan di [`01-arsitektur.md`](01-arsitektur.md)).

_[isi sendiri: tambahkan bagian yang kamu tulis/ubah sendiri di luar bantuan AI,
atau bagian yang kamu pelajari ulang sampai paham.]_
