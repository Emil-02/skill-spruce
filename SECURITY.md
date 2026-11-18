# Security Policy

`skill-spruce` hanya membaca file skill dan menulis laporan. Ia tidak menjalankan perintah yang ditemukan di dalam skill, tidak mengirim data ke layanan eksternal, dan tidak mengubah file yang diaudit. Pola berbahaya seperti `curl | bash` atau `auto-push` hanya dilaporkan.

## Melaporkan kerentanan

Sebaiknya jangan buka issue publik. Kirim laporan lewat [GitHub Security Advisories](https://github.com/Emil-02/skill-spruce/security/advisories/new) supaya detailnya tetap privat. Saya biasanya merespons dalam 7 hari.

## Dukungan

Rilis `0.1.x` masih saya dukung dengan perbaikan keamanan.
