# Rencana skill-spruce

## Mengapa proyek ini

Tren GitHub 2026: repo yang tumbuh paling cepat bukan agent, melainkan **skill dan konfigurasi untuk agent**. Semua orang membuat `SKILL.md`, hampir tidak ada yang memeriksa apakah instruksi itu bagus. `skill-spruce` mengisi celah itu.

Saya sudah punya resep yang terbukti lewat `gh-spruce`: pindai → skor 0–100 → jelaskan masalah → kasih badge → perbaikan otomatis. Resep itu saya pakai ulang untuk domain skill.

## Prinsip

- **Read-only.** Tidak mengubah file skill tanpa konfirmasi.
- **Menjelaskan, bukan menghakimi.** Setiap temuan punya alasan dan saran.
- **Ringan.** Satu binary Go, tanpa dependensi eksternal, tanpa jaringan.
- **Bisa dipakai di CI.** Keluaran JSON/SARIF dan mode `--ci`.

## Tujuh aturan

| ID | Bobot | Severity | Gagal jika |
|---|---|---|---|
| `frontmatter` | 20 | high | Tidak ada blok `---`, atau `name`/`description` kosong |
| `description` | 20 | high | Pendek, generik, atau tidak menyebut kapan skill dipakai |
| `dangerous` | 20 | critical | Ada `curl \| bash`, `rm -rf /`, `auto-push`, atau perintah print secret |
| `hard-rules` | 15 | high | Tidak ada batas eksplisit seperti "jangan", "dilarang", "never" |
| `length` | 10 | medium | Isi lebih dari 8.000 karakter tanpa `references/` |
| `agent-coverage` | 10 | low | Tidak ada penanda agent lain selain satu format |
| `progressive` | 5 | low | Tidak ada tautan ke `references/` padahal instruksi panjang |

## Fase

### Fase 0 — Fondasi (selesai)
- Nama, lisensi MIT, `CODEOWNERS`, CI, README, deskripsi & topics repo.
- Dokumen rencana ini.

### Fase 1 — Mesin inti (selesai)
- Pemindai file skill di seluruh direktori (abaikan `node_modules`, `.git`, `dist`, `vendor`).
- Tujuh aturan di atas.
- Perhitungan skor: `earned / total × 100`, `warning` dapat setengah bobot, temuan `critical` membatasi skor ke 40.
- Keluaran `text`, `json`, `sarif`, `badge`.

### Fase 2 — Pengalaman pakai
- `skill-spruce fix --dry-run` — usulan perbaikan otomatis untuk `frontmatter` dan `description`.
- Konfigurasi `.skillspruce.yml`: matikan aturan, ubah bobot, atur ambang.
- Bandingkan dua versi skill (`skill-spruce diff`) untuk melihat regresi kualitas.

### Fase 3 — Ekosistem
- GitHub Action resmi (`action.yml`).
- Badge permanen lewat endpoint `/api/badge?repo=...` (mengikuti pola `gh-spruce`).
- Dukungan format skill Cursor (`.mdc`) dan `CLAUDE.md` lebih dalam.

### Fase 4 — Komunitas
- `good first issue`: aturan baru, deteksi agent baru, terjemahan.
- Katalog contoh `SKILL.md` baik dan buruk di `examples/`.
- Artikel: "Saya mengaudit skill agent saya sendiri dan dapat 88/100".

## Metrik sukses

- Bisa dijalankan di repo sendiri (`skill-spruce check .` pada `jaga-skill`).
- Skor gh-spruce repo ini tetap 100/100.
- Test Go lolos, binary bisa di-`go install`.
- Orang lain bisa menambah aturan baru tanpa mengubah inti (interface `rules.Rule`).

## Atribusi

Repository ini saya kelola dari akun GitHub `Emil-02`. `CODEOWNERS` dan identitas commit mengacu ke akun tersebut.
