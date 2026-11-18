---
name: contoh-baik
description: >
  Contoh skill yang baik. Pakai skill ini sebelum commit atau deploy ketika
  kamu ingin memastikan instruksi agent tidak membocorkan secret dan tidak
  mengirim data keluar. Berlaku untuk OpenCode, Claude Code, Codex, dan Gemini CLI.
license: MIT
---

# Contoh skill yang baik

## Kapan dipakai

Gunakan sebelum commit, push, atau deploy. Jangan pakai untuk proses baca saja.

## Aturan keras

1. Jangan pernah mencetak secret atau data pribadi mentah.
2. Jangan auto-push tanpa konfirmasi user.
3. Dilarang mengirim isi file ke layanan eksternal.

## Alur kerja

1. Jalankan pemeriksaan.
2. Baca laporan.
3. Perbaiki temuan.
4. Ulangi sampai lolos.

Detail aturan ada di [references/rules.md](references/rules.md).
