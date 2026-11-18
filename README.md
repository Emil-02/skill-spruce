# skill-spruce

[![CI](https://github.com/Emil-02/skill-spruce/actions/workflows/ci.yml/badge.svg)](https://github.com/Emil-02/skill-spruce/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Saya bikin `skill-spruce` karena skill AI makin banyak, tapi kualitasnya belum ada yang mengawasi. Kalau `gh-spruce` menilai kesehatan repo, `skill-spruce` menilai kesehatan **instruksi agent** — file `SKILL.md`, `AGENTS.md`, `GEMINI.md`, dan `.cursor/rules/*.mdc`.

```sh
go install github.com/Emil-02/skill-spruce/cmd/skill-spruce@latest
skill-spruce check .
```

## Yang diperiksa

| Aturan | Bobot | Pertanyaan yang dijawab |
|---|---|---|
| `frontmatter` | 20 | Apakah `SKILL.md` punya `name` dan `description`? |
| `description` | 20 | Apakah deskripsinya spesifik dan punya pemicu pakai? |
| `dangerous` | 20 | Apakah ada instruksi berbahaya (`curl \| bash`, auto-push, secret mentah)? |
| `hard-rules` | 15 | Apakah ada batas eksplisit untuk agent? |
| `length` | 10 | Apakah instruksi terlalu panjang dan boros konteks? |
| `agent-coverage` | 10 | Apakah skill bisa dipakai lebih dari satu agent? |
| `progressive` | 5 | Apakah detail dipindah ke `references/`? |

Skor 0–100, grade A+ sampai F. Temuan `critical` membatasi skor maksimal 40 — sama seperti `gh-spruce`.

## Pemakaian

```sh
skill-spruce check .                       # ringkas di terminal
skill-spruce check . --format json
skill-spruce check . --format sarif -o skill.sarif
skill-spruce check . --format badge -o badge.svg
skill-spruce check . --ci --fail-on high   # untuk pipeline
skill-spruce rules                         # lihat daftar aturan
```

Kode keluar `0` lolos, `1` ada temuan di atas ambang, `2` error.

## Contoh hasil

```
skill-spruce — /home/me/project
────────────────────────────────────────────
✓ frontmatter     HIGH
! description     HIGH
    SKILL.md:1 — deskripsi terlalu umum; sebutkan kapan skill dipakai
✓ dangerous       CRITICAL
✓ hard-rules      HIGH
✓ length          MEDIUM
✓ agent-coverage  LOW
✓ progressive     LOW
────────────────────────────────────────────
Skor: 88/100 (A)
```

## Rencana

Rencana lengkap dan alasan tiap fase ada di [`docs/PLAN.md`](docs/PLAN.md). Arsitektur kode ada di [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Hubungan dengan project saya

- [`jaga`](https://github.com/Emil-02/jaga) menjaga **kode**.
- [`jaga-skill`](https://github.com/Emil-02/jaga-skill) menjaga **cara agent bekerja**.
- `skill-spruce` menjaga **mutu aturan agent itu sendiri**.

## Lisensi

MIT. Dibuat dan dipelihara oleh [@Emil-02](https://github.com/Emil-02).
